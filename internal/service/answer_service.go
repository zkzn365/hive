package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"answer/internal/base/constant"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service/activity"
	"answer/internal/service/activity_common"
	"answer/internal/service/activity_queue"
	answercommon "answer/internal/service/answer_common"
	collectioncommon "answer/internal/service/collection_common"
	"answer/internal/service/notice_queue"
	"answer/internal/service/permission"
	questioncommon "answer/internal/service/question_common"
	"answer/internal/service/revision_common"
	usercommon "answer/internal/service/user_common"

	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

// AnswerService user service
type AnswerService struct {
	answerRepo            answercommon.AnswerRepo
	questionRepo          questioncommon.QuestionRepo
	questionCommon        *questioncommon.QuestionCommon
	answerActivityService *activity.AnswerActivityService
	userCommon            *usercommon.UserCommon
	collectionCommon      *collectioncommon.CollectionCommon
	userRepo              usercommon.UserRepo
	revisionService       *revision_common.RevisionService
	AnswerCommon          *answercommon.AnswerCommon
	voteRepo              activity_common.VoteRepo
}

func NewAnswerService(
	answerRepo answercommon.AnswerRepo,
	questionRepo questioncommon.QuestionRepo,
	questionCommon *questioncommon.QuestionCommon,
	userCommon *usercommon.UserCommon,
	collectionCommon *collectioncommon.CollectionCommon,
	userRepo usercommon.UserRepo,
	revisionService *revision_common.RevisionService,
	answerAcceptActivityRepo *activity.AnswerActivityService,
	answerCommon *answercommon.AnswerCommon,
	voteRepo activity_common.VoteRepo,
) *AnswerService {
	return &AnswerService{
		answerRepo:            answerRepo,
		questionRepo:          questionRepo,
		userCommon:            userCommon,
		collectionCommon:      collectionCommon,
		questionCommon:        questionCommon,
		userRepo:              userRepo,
		revisionService:       revisionService,
		answerActivityService: answerAcceptActivityRepo,
		AnswerCommon:          answerCommon,
		voteRepo:              voteRepo,
	}
}

// RemoveAnswer delete answer
func (as *AnswerService) RemoveAnswer(ctx context.Context, req *schema.RemoveAnswerReq) (err error) {
	answerInfo, exist, err := as.answerRepo.GetByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	if !req.IsAdmin {
		if answerInfo.UserID != req.UserID {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		if answerInfo.VoteCount > 0 {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		if answerInfo.Adopted == schema.AnswerAdoptedEnable {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		questionInfo, exist, err := as.questionRepo.GetQuestion(ctx, answerInfo.QuestionID)
		if err != nil {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		if !exist {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		if questionInfo.AnswerCount > 1 {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
		if questionInfo.AcceptedAnswerID != "" {
			return errors.BadRequest(reason.AnswerCannotDeleted)
		}
	}

	// user add question count
	err = as.questionCommon.UpdateAnswerCount(ctx, answerInfo.QuestionID, -1)
	if err != nil {
		log.Error("IncreaseAnswerCount error", err.Error())
	}

	err = as.userCommon.UpdateAnswerCount(ctx, answerInfo.UserID, -1)
	if err != nil {
		log.Error("user IncreaseAnswerCount error", err.Error())
	}

	err = as.answerRepo.RemoveAnswer(ctx, req.ID)
	if err != nil {
		return err
	}
	err = as.answerActivityService.DeleteAnswer(ctx, answerInfo.ID, answerInfo.CreatedAt, answerInfo.VoteCount)
	if err != nil {
		log.Errorf("delete answer activity change failed: %s", err.Error())
	}
	activity_queue.AddActivity(&schema.ActivityMsg{
		UserID:           req.UserID,
		ObjectID:         answerInfo.ID,
		OriginalObjectID: answerInfo.ID,
		ActivityTypeKey:  constant.ActAnswerDeleted,
	})
	return
}

func (as *AnswerService) Insert(ctx context.Context, req *schema.AnswerAddReq) (string, error) {
	questionInfo, exist, err := as.questionRepo.GetQuestion(ctx, req.QuestionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.BadRequest(reason.QuestionNotFound)
	}
	insertData := new(entity.Answer)
	insertData.UserID = req.UserID
	insertData.OriginalText = req.Content
	insertData.ParsedText = req.HTML
	insertData.Adopted = schema.AnswerAdoptedFailed
	insertData.QuestionID = req.QuestionID
	insertData.RevisionID = "0"
	insertData.LastEditUserID = "0"
	insertData.Status = entity.AnswerStatusAvailable
	//insertData.UpdatedAt = now
	if err = as.answerRepo.AddAnswer(ctx, insertData); err != nil {
		return "", err
	}
	err = as.questionCommon.UpdateAnswerCount(ctx, req.QuestionID, 1)
	if err != nil {
		log.Error("IncreaseAnswerCount error", err.Error())
	}
	err = as.questionCommon.UpdateLastAnswer(ctx, req.QuestionID, insertData.ID)
	if err != nil {
		log.Error("UpdateLastAnswer error", err.Error())
	}
	err = as.questionCommon.UpdataPostTime(ctx, req.QuestionID)
	if err != nil {
		return insertData.ID, err
	}

	err = as.userCommon.UpdateAnswerCount(ctx, req.UserID, 1)
	if err != nil {
		log.Error("user IncreaseAnswerCount error", err.Error())
	}

	revisionDTO := &schema.AddRevisionDTO{
		UserID:   insertData.UserID,
		ObjectID: insertData.ID,
		Title:    "",
	}
	infoJSON, _ := json.Marshal(insertData)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := as.revisionService.AddRevision(ctx, revisionDTO, true)
	if err != nil {
		return insertData.ID, err
	}
	as.notificationAnswerTheQuestion(ctx, questionInfo.UserID, insertData.ID, req.UserID)

	activity_queue.AddActivity(&schema.ActivityMsg{
		UserID:           insertData.UserID,
		ObjectID:         insertData.ID,
		OriginalObjectID: insertData.ID,
		ActivityTypeKey:  constant.ActAnswerAnswered,
		RevisionID:       revisionID,
	})
	activity_queue.AddActivity(&schema.ActivityMsg{
		UserID:           insertData.UserID,
		ObjectID:         insertData.ID,
		OriginalObjectID: questionInfo.ID,
		ActivityTypeKey:  constant.ActQuestionAnswered,
	})
	return insertData.ID, nil
}

func (as *AnswerService) Update(ctx context.Context, req *schema.AnswerUpdateReq) (string, error) {
	//req.NoNeedReview //true 不需要审核
	var canUpdate bool
	_, existUnreviewed, err := as.revisionService.ExistUnreviewedByObjectID(ctx, req.ID)
	if err != nil {
		return "", err
	}
	if existUnreviewed {
		err = errors.BadRequest(reason.AnswerCannotUpdate)
		return "", err
	}

	questionInfo, exist, err := as.questionRepo.GetQuestion(ctx, req.QuestionID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", errors.BadRequest(reason.QuestionNotFound)
	}

	answerInfo, exist, err := as.answerRepo.GetByID(ctx, req.ID)
	if err != nil {
		return "", err
	}
	if !exist {
		return "", nil
	}

	//If the content is the same, ignore it
	if answerInfo.OriginalText == req.Content {
		return "", nil
	}

	now := time.Now()
	insertData := new(entity.Answer)
	insertData.ID = req.ID
	insertData.UserID = answerInfo.UserID
	insertData.QuestionID = req.QuestionID
	insertData.OriginalText = req.Content
	insertData.ParsedText = req.HTML
	insertData.UpdatedAt = now

	insertData.LastEditUserID = "0"
	if answerInfo.UserID != req.UserID {
		insertData.LastEditUserID = req.UserID
	}

	revisionDTO := &schema.AddRevisionDTO{
		UserID:   req.UserID,
		ObjectID: req.ID,
		Title:    "",
		Log:      req.EditSummary,
	}

	if req.NoNeedReview || answerInfo.UserID == req.UserID {
		canUpdate = true
	}

	if !canUpdate {
		revisionDTO.Status = entity.RevisionUnreviewedStatus
	} else {
		if err = as.answerRepo.UpdateAnswer(ctx, insertData, []string{"original_text", "parsed_text", "updated_at", "last_edit_user_id"}); err != nil {
			return "", err
		}
		err = as.questionCommon.UpdataPostTime(ctx, req.QuestionID)
		if err != nil {
			return insertData.ID, err
		}
		as.notificationUpdateAnswer(ctx, questionInfo.UserID, insertData.ID, req.UserID)
		revisionDTO.Status = entity.RevisionReviewPassStatus
	}

	infoJSON, _ := json.Marshal(insertData)
	revisionDTO.Content = string(infoJSON)
	revisionID, err := as.revisionService.AddRevision(ctx, revisionDTO, true)
	if err != nil {
		return insertData.ID, err
	}
	if canUpdate {
		activity_queue.AddActivity(&schema.ActivityMsg{
			UserID:           insertData.UserID,
			ObjectID:         insertData.ID,
			OriginalObjectID: insertData.ID,
			ActivityTypeKey:  constant.ActAnswerEdited,
			RevisionID:       revisionID,
		})
	}

	return insertData.ID, nil
}

// UpdateAdopted
func (as *AnswerService) UpdateAdopted(ctx context.Context, req *schema.AnswerAdoptedReq) error {
	if req.AnswerID == "" {
		req.AnswerID = "0"
	}
	if req.UserID == "" {
		return nil
	}

	newAnswerInfo := &entity.Answer{}
	newAnswerInfoexist := false
	var err error

	if req.AnswerID != "0" {
		newAnswerInfo, newAnswerInfoexist, err = as.answerRepo.GetByID(ctx, req.AnswerID)
		if err != nil {
			return err
		}
		if !newAnswerInfoexist {
			return errors.BadRequest(reason.AnswerNotFound)
		}
	}

	questionInfo, exist, err := as.questionRepo.GetQuestion(ctx, req.QuestionID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.BadRequest(reason.QuestionNotFound)
	}
	if questionInfo.UserID != req.UserID {
		return fmt.Errorf("no permission to set answer")
	}
	if questionInfo.AcceptedAnswerID == req.AnswerID {
		return nil
	}

	var oldAnswerInfo *entity.Answer
	if len(questionInfo.AcceptedAnswerID) > 0 && questionInfo.AcceptedAnswerID != "0" {
		oldAnswerInfo, _, err = as.answerRepo.GetByID(ctx, questionInfo.AcceptedAnswerID)
		if err != nil {
			return err
		}
	}

	err = as.answerRepo.UpdateAdopted(ctx, req.AnswerID, req.QuestionID)
	if err != nil {
		return err
	}

	err = as.questionCommon.UpdateAccepted(ctx, req.QuestionID, req.AnswerID)
	if err != nil {
		log.Error("UpdateLastAnswer error", err.Error())
	}

	as.updateAnswerRank(ctx, req.UserID, questionInfo, newAnswerInfo, oldAnswerInfo)
	return nil
}

func (as *AnswerService) updateAnswerRank(ctx context.Context, userID string,
	questionInfo *entity.Question, newAnswerInfo *entity.Answer, oldAnswerInfo *entity.Answer,
) {
	// if this question is already been answered, should cancel old answer rank
	if oldAnswerInfo != nil {
		err := as.answerActivityService.CancelAcceptAnswer(
			ctx, questionInfo.AcceptedAnswerID, questionInfo.ID, questionInfo.UserID, oldAnswerInfo.UserID)
		if err != nil {
			log.Error(err)
		}
	}
	if newAnswerInfo.ID != "" {
		err := as.answerActivityService.AcceptAnswer(
			ctx, newAnswerInfo.ID, questionInfo.ID, questionInfo.UserID, newAnswerInfo.UserID, newAnswerInfo.UserID == userID)
		if err != nil {
			log.Error(err)
		}
	}
}

func (as *AnswerService) Get(ctx context.Context, answerID, loginUserID string) (*schema.AnswerInfo, *schema.QuestionInfo, bool, error) {
	answerInfo, has, err := as.answerRepo.GetByID(ctx, answerID)
	if err != nil {
		return nil, nil, has, err
	}
	info := as.ShowFormat(ctx, answerInfo)
	// todo questionFunc
	questionInfo, err := as.questionCommon.Info(ctx, answerInfo.QuestionID, loginUserID)
	if err != nil {
		return nil, nil, has, err
	}
	// todo UserFunc

	userIds := make([]string, 0)
	userIds = append(userIds, answerInfo.UserID)
	userIds = append(userIds, answerInfo.LastEditUserID)
	userInfoMap, err := as.userCommon.BatchUserBasicInfoByID(ctx, userIds)
	if err != nil {
		return nil, nil, has, err
	}

	_, ok := userInfoMap[answerInfo.UserID]
	if ok {
		info.UserInfo = userInfoMap[answerInfo.UserID]
	}
	_, ok = userInfoMap[answerInfo.LastEditUserID]
	if ok {
		info.UpdateUserInfo = userInfoMap[answerInfo.LastEditUserID]
	}

	if loginUserID == "" {
		return info, questionInfo, has, nil
	}

	info.VoteStatus = as.voteRepo.GetVoteStatus(ctx, answerID, loginUserID)

	CollectedMap, err := as.collectionCommon.SearchObjectCollected(ctx, loginUserID, []string{answerInfo.ID})
	if err != nil {
		log.Error("CollectionFunc.SearchObjectCollected error", err)
	}
	_, ok = CollectedMap[answerInfo.ID]
	if ok {
		info.Collected = true
	}

	return info, questionInfo, has, nil
}

func (as *AnswerService) AdminSetAnswerStatus(ctx context.Context, req *schema.AdminSetAnswerStatusRequest) error {
	setStatus, ok := entity.CmsAnswerSearchStatus[req.StatusStr]
	if !ok {
		return fmt.Errorf("question status does not exist")
	}
	answerInfo, exist, err := as.answerRepo.GetAnswer(ctx, req.AnswerID)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("answer does not exist")
	}
	answerInfo.Status = setStatus
	err = as.answerRepo.UpdateAnswerStatus(ctx, answerInfo)
	if err != nil {
		return err
	}

	if setStatus == entity.AnswerStatusDeleted {
		err = as.answerActivityService.DeleteAnswer(ctx, answerInfo.ID, answerInfo.CreatedAt, answerInfo.VoteCount)
		if err != nil {
			log.Errorf("admin delete question then rank rollback error %s", err.Error())
		} else {
			activity_queue.AddActivity(&schema.ActivityMsg{
				UserID:           req.UserID,
				ObjectID:         answerInfo.ID,
				OriginalObjectID: answerInfo.ID,
				ActivityTypeKey:  constant.ActAnswerDeleted,
			})
		}
	}

	msg := &schema.NotificationMsg{}
	msg.ObjectID = answerInfo.ID
	msg.Type = schema.NotificationTypeInbox
	msg.ReceiverUserID = answerInfo.UserID
	msg.TriggerUserID = answerInfo.UserID
	msg.ObjectType = constant.AnswerObjectType
	msg.NotificationAction = constant.YourAnswerWasDeleted
	notice_queue.AddNotification(msg)

	return nil
}

func (as *AnswerService) SearchList(ctx context.Context, req *schema.AnswerListReq) ([]*schema.AnswerInfo, int64, error) {
	list := make([]*schema.AnswerInfo, 0)
	dbSearch := entity.AnswerSearch{}
	dbSearch.QuestionID = req.QuestionID
	dbSearch.Page = req.Page
	dbSearch.PageSize = req.PageSize
	dbSearch.Order = req.Order
	answerOriginalList, count, err := as.answerRepo.SearchList(ctx, &dbSearch)
	if err != nil {
		return list, count, err
	}
	answerList, err := as.SearchFormatInfo(ctx, answerOriginalList, req)
	if err != nil {
		return answerList, count, err
	}
	return answerList, count, nil
}

func (as *AnswerService) SearchFormatInfo(ctx context.Context, answers []*entity.Answer, req *schema.AnswerListReq) (
	[]*schema.AnswerInfo, error) {
	list := make([]*schema.AnswerInfo, 0)
	objectIDs := make([]string, 0)
	userIDs := make([]string, 0)
	for _, info := range answers {
		item := as.ShowFormat(ctx, info)
		list = append(list, item)
		objectIDs = append(objectIDs, info.ID)
		userIDs = append(userIDs, info.UserID)
		userIDs = append(userIDs, info.LastEditUserID)
		if req.UserID != "" {
			item.VoteStatus = as.voteRepo.GetVoteStatus(ctx, item.ID, req.UserID)
		}
	}
	userInfoMap, err := as.userCommon.BatchUserBasicInfoByID(ctx, userIDs)
	if err != nil {
		return list, err
	}
	for _, item := range list {
		_, ok := userInfoMap[item.UserID]
		if ok {
			item.UserInfo = userInfoMap[item.UserID]
		}
		_, ok = userInfoMap[item.UpdateUserID]
		if ok {
			item.UpdateUserInfo = userInfoMap[item.UpdateUserID]
		}
	}

	if req.UserID == "" {
		return list, nil
	}

	searchObjectCollected, err := as.collectionCommon.SearchObjectCollected(ctx, req.UserID, objectIDs)
	if err != nil {
		return nil, err
	}

	for _, item := range list {
		_, ok := searchObjectCollected[item.ID]
		if ok {
			item.Collected = true
		}
	}

	for _, item := range list {
		item.MemberActions = permission.GetAnswerPermission(ctx, req.UserID, item.UserID, req.CanEdit, req.CanDelete)
	}
	return list, nil
}

func (as *AnswerService) ShowFormat(ctx context.Context, data *entity.Answer) *schema.AnswerInfo {
	return as.AnswerCommon.ShowFormat(ctx, data)
}

func (as *AnswerService) notificationUpdateAnswer(ctx context.Context, questionUserID, answerID, answerUserID string) {
	msg := &schema.NotificationMsg{
		TriggerUserID:  answerUserID,
		ReceiverUserID: questionUserID,
		Type:           schema.NotificationTypeInbox,
		ObjectID:       answerID,
	}
	msg.ObjectType = constant.AnswerObjectType
	msg.NotificationAction = constant.UpdateAnswer
	notice_queue.AddNotification(msg)
}

func (as *AnswerService) notificationAnswerTheQuestion(ctx context.Context, questionUserID, answerID, answerUserID string) {
	msg := &schema.NotificationMsg{
		TriggerUserID:  answerUserID,
		ReceiverUserID: questionUserID,
		Type:           schema.NotificationTypeInbox,
		ObjectID:       answerID,
	}
	msg.ObjectType = constant.AnswerObjectType
	msg.NotificationAction = constant.AnswerTheQuestion
	notice_queue.AddNotification(msg)
}
