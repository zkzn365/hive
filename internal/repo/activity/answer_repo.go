package activity

import (
	"context"
	"time"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service/activity"
	"answer/internal/service/activity_common"
	"answer/internal/service/notice_queue"
	"answer/internal/service/rank"
	"answer/pkg/converter"

	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"xorm.io/xorm"
)

// AnswerActivityRepo answer accepted
type AnswerActivityRepo struct {
	data         *data.Data
	activityRepo activity_common.ActivityRepo
	userRankRepo rank.UserRankRepo
}

const (
	acceptAction         = "accept"
	acceptedAction       = "accepted"
	acceptCancelAction   = "accept_cancel"
	acceptedCancelAction = "accepted_cancel"
)

var (
	acceptActionList       = []string{acceptAction, acceptedAction}
	acceptCancelActionList = []string{acceptCancelAction, acceptedCancelAction}
)

// NewAnswerActivityRepo new repository
func NewAnswerActivityRepo(
	data *data.Data,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
) activity.AnswerActivityRepo {
	return &AnswerActivityRepo{
		data:         data,
		activityRepo: activityRepo,
		userRankRepo: userRankRepo,
	}
}

// NewQuestionActivityRepo new repository
func NewQuestionActivityRepo(
	data *data.Data,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
) activity.QuestionActivityRepo {
	return &AnswerActivityRepo{
		data:         data,
		activityRepo: activityRepo,
		userRankRepo: userRankRepo,
	}
}

func (ar *AnswerActivityRepo) DeleteQuestion(ctx context.Context, questionID string) (err error) {
	questionInfo := &entity.Question{}
	exist, err := ar.data.DB.Where("id = ?", questionID).Get(questionInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil
	}

	// get all this object activity
	activityList := make([]*entity.Activity, 0)
	session := ar.data.DB.Where("has_rank = 1")
	session.Where("cancelled = ?", entity.ActivityAvailable)
	err = session.Find(&activityList, &entity.Activity{ObjectID: questionID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(activityList) == 0 {
		return nil
	}

	log.Infof("questionInfo %s deleted will rollback activity %d", questionID, len(activityList))

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		for _, act := range activityList {
			log.Infof("user %s rollback rank %d", act.UserID, -act.Rank)
			_, e := ar.userRankRepo.TriggerUserRank(
				ctx, session, act.UserID, -act.Rank, act.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}

			if _, e := session.Where("id = ?", act.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{Cancelled: entity.ActivityCancelled, CancelledAt: time.Now()}); e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	// get all answers
	answerList := make([]*entity.Answer, 0)
	err = ar.data.DB.Find(&answerList, &entity.Answer{QuestionID: questionID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, answerInfo := range answerList {
		err = ar.DeleteAnswer(ctx, answerInfo.ID)
		if err != nil {
			log.Error(err)
		}
	}
	return
}

// AcceptAnswer accept other answer
func (ar *AnswerActivityRepo) AcceptAnswer(ctx context.Context,
	answerObjID, questionObjID, questionUserID, answerUserID string, isSelf bool,
) (err error) {
	addActivityList := make([]*entity.Activity, 0)
	for _, action := range acceptActionList {
		// get accept answer need add rank amount
		activityType, deltaRank, hasRank, e := ar.activityRepo.GetActivityTypeByObjID(ctx, answerObjID, action)
		if e != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
		}
		addActivity := &entity.Activity{
			ObjectID:         answerObjID,
			OriginalObjectID: questionObjID,
			ActivityType:     activityType,
			Rank:             deltaRank,
			HasRank:          hasRank,
		}
		if action == acceptAction {
			addActivity.UserID = questionUserID
			addActivity.TriggerUserID = converter.StringToInt64(answerUserID)
			addActivity.OriginalObjectID = questionObjID // if activity is 'accept' means this question is accept the answer.
		} else {
			addActivity.UserID = answerUserID
			addActivity.TriggerUserID = converter.StringToInt64(answerUserID)
			addActivity.OriginalObjectID = answerObjID // if activity is 'accepted' means this answer was accepted.
		}
		if isSelf {
			addActivity.Rank = 0
		}
		addActivityList = append(addActivityList, addActivity)
	}

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		for _, addActivity := range addActivityList {
			existsActivity, exists, e := ar.activityRepo.GetActivity(
				ctx, session, answerObjID, addActivity.UserID, addActivity.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
			if exists && existsActivity.Cancelled == entity.ActivityAvailable {
				continue
			}

			reachStandard, e := ar.userRankRepo.TriggerUserRank(
				ctx, session, addActivity.UserID, addActivity.Rank, addActivity.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
			if reachStandard {
				addActivity.Rank = 0
			}

			if exists {
				if _, e = session.Where("id = ?", existsActivity.ID).Cols("`cancelled`").
					Update(&entity.Activity{Cancelled: entity.ActivityAvailable}); e != nil {
					return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
				}
			} else {
				if _, e = session.Insert(addActivity); e != nil {
					return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
				}
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	for _, act := range addActivityList {
		msg := &schema.NotificationMsg{
			Type:           schema.NotificationTypeAchievement,
			ObjectID:       act.ObjectID,
			ReceiverUserID: act.UserID,
		}
		if act.UserID == questionUserID {
			msg.TriggerUserID = answerUserID
			msg.ObjectType = constant.AnswerObjectType
		} else {
			msg.TriggerUserID = questionUserID
			msg.ObjectType = constant.AnswerObjectType
		}
		notice_queue.AddNotification(msg)
	}

	for _, act := range addActivityList {
		msg := &schema.NotificationMsg{
			ReceiverUserID: act.UserID,
			Type:           schema.NotificationTypeInbox,
			ObjectID:       act.ObjectID,
		}
		if act.UserID != questionUserID {
			msg.TriggerUserID = questionUserID
			msg.ObjectType = constant.AnswerObjectType
			msg.NotificationAction = constant.AdoptAnswer
			notice_queue.AddNotification(msg)
		}
	}
	return err
}

// CancelAcceptAnswer accept other answer
func (ar *AnswerActivityRepo) CancelAcceptAnswer(ctx context.Context,
	answerObjID, questionObjID, questionUserID, answerUserID string,
) (err error) {
	addActivityList := make([]*entity.Activity, 0)
	for _, action := range acceptActionList {
		// get accept answer need add rank amount
		activityType, deltaRank, hasRank, e := ar.activityRepo.GetActivityTypeByObjID(ctx, answerObjID, action)
		if e != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
		}
		addActivity := &entity.Activity{
			ObjectID:     answerObjID,
			ActivityType: activityType,
			Rank:         -deltaRank,
			HasRank:      hasRank,
		}
		if action == acceptAction {
			addActivity.UserID = questionUserID
			addActivity.OriginalObjectID = questionObjID
		} else {
			addActivity.UserID = answerUserID
			addActivity.OriginalObjectID = answerObjID
		}
		addActivityList = append(addActivityList, addActivity)
	}

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		for _, addActivity := range addActivityList {
			existsActivity, exists, e := ar.activityRepo.GetActivity(
				ctx, session, answerObjID, addActivity.UserID, addActivity.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
			if exists && existsActivity.Cancelled == entity.ActivityCancelled {
				continue
			}
			if !exists {
				continue
			}

			_, e = ar.userRankRepo.TriggerUserRank(
				ctx, session, addActivity.UserID, addActivity.Rank, addActivity.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}

			if _, e := session.Where("id = ?", existsActivity.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{Cancelled: entity.ActivityCancelled, CancelledAt: time.Now()}); e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	for _, act := range addActivityList {
		msg := &schema.NotificationMsg{
			ReceiverUserID: act.UserID,
			Type:           schema.NotificationTypeAchievement,
			ObjectID:       act.ObjectID,
		}
		if act.UserID == questionUserID {
			msg.TriggerUserID = answerUserID
			msg.ObjectType = constant.QuestionObjectType
		} else {
			msg.TriggerUserID = questionUserID
			msg.ObjectType = constant.AnswerObjectType
		}
		notice_queue.AddNotification(msg)
	}
	return err
}

func (ar *AnswerActivityRepo) DeleteAnswer(ctx context.Context, answerID string) (err error) {
	answerInfo := &entity.Answer{}
	exist, err := ar.data.DB.Where("id = ?", answerID).Get(answerInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil
	}

	// get all this object activity
	activityList := make([]*entity.Activity, 0)
	session := ar.data.DB.Where("has_rank = 1")
	session.Where("cancelled = ?", entity.ActivityAvailable)
	err = session.Find(&activityList, &entity.Activity{ObjectID: answerID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(activityList) == 0 {
		return nil
	}

	log.Infof("answerInfo %s deleted will rollback activity %d", answerID, len(activityList))

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		for _, act := range activityList {
			log.Infof("user %s rollback rank %d", act.UserID, -act.Rank)
			_, e := ar.userRankRepo.TriggerUserRank(
				ctx, session, act.UserID, -act.Rank, act.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}

			if _, e := session.Where("id = ?", act.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{Cancelled: entity.ActivityCancelled, CancelledAt: time.Now()}); e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return
}
