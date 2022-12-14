package activity

import (
	"context"
	"strings"
	"time"

	"answer/pkg/converter"

	"answer/internal/base/pager"
	"answer/internal/service/config"
	"answer/internal/service/notice_queue"
	"answer/internal/service/rank"
	"answer/pkg/obj"

	"xorm.io/builder"

	"answer/internal/service/activity_common"
	"answer/internal/service/unique"

	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service"

	"github.com/segmentfault/pacman/errors"
	"xorm.io/xorm"
)

// VoteRepo activity repository
type VoteRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
	configRepo   config.ConfigRepo
	activityRepo activity_common.ActivityRepo
	userRankRepo rank.UserRankRepo
	voteCommon   activity_common.VoteRepo
}

// NewVoteRepo new repository
func NewVoteRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
	configRepo config.ConfigRepo,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
	voteCommon activity_common.VoteRepo,
) service.VoteRepo {
	return &VoteRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
		configRepo:   configRepo,
		activityRepo: activityRepo,
		userRankRepo: userRankRepo,
		voteCommon:   voteCommon,
	}
}

var LimitUpActions = map[string][]string{
	"question": {"vote_up", "voted_up"},
	"answer":   {"vote_up", "voted_up"},
	"comment":  {"vote_up"},
}

var LimitDownActions = map[string][]string{
	"question": {"vote_down", "voted_down"},
	"answer":   {"vote_down", "voted_down"},
	"comment":  {"vote_down"},
}

func (vr *VoteRepo) vote(ctx context.Context, objectID string, userID, objectUserID string, actions []string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		result = nil
		for _, action := range actions {
			var (
				existsActivity entity.Activity
				insertActivity entity.Activity
				has            bool
				triggerUserID,
				activityUserID string
				activityType, deltaRank, hasRank int
			)

			activityUserID, activityType, deltaRank, hasRank, err = vr.CheckRank(ctx, objectID, objectUserID, userID, action)
			if err != nil {
				return
			}

			triggerUserID = userID
			if userID == activityUserID {
				triggerUserID = "0"
			}

			// check is voted up
			has, _ = session.
				Where(builder.Eq{"object_id": objectID}).
				And(builder.Eq{"user_id": activityUserID}).
				And(builder.Eq{"trigger_user_id": triggerUserID}).
				And(builder.Eq{"activity_type": activityType}).
				Get(&existsActivity)

			// is is voted,return
			if has && existsActivity.Cancelled == entity.ActivityAvailable {
				return
			}

			insertActivity = entity.Activity{
				ObjectID:         objectID,
				OriginalObjectID: objectID,
				UserID:           activityUserID,
				TriggerUserID:    converter.StringToInt64(triggerUserID),
				ActivityType:     activityType,
				Rank:             deltaRank,
				HasRank:          hasRank,
				Cancelled:        entity.ActivityAvailable,
			}

			// trigger user rank and send notification
			if hasRank != 0 {
				var isReachStandard bool
				isReachStandard, err = vr.userRankRepo.TriggerUserRank(ctx, session, activityUserID, deltaRank, activityType)
				if err != nil {
					return nil, err
				}
				if isReachStandard {
					insertActivity.Rank = 0
				}

				vr.sendNotification(ctx, activityUserID, objectUserID, objectID)
			}

			if has {
				if _, err = session.Where("id = ?", existsActivity.ID).Cols("`cancelled`").
					Update(&entity.Activity{
						Cancelled: entity.ActivityAvailable,
					}); err != nil {
					return
				}
			} else {
				_, err = session.Insert(&insertActivity)
				if err != nil {
					return nil, err
				}
			}

			// update votes
			if action == "vote_down" || action == "vote_up" {
				votes := 1
				if action == "vote_down" {
					votes = -1
				}
				err = vr.updateVotes(ctx, session, objectID, votes)
				if err != nil {
					return
				}
			}
		}
		return
	})
	if err != nil {
		return
	}

	resp, err = vr.GetVoteResultByObjectId(ctx, objectID)
	resp.VoteStatus = vr.voteCommon.GetVoteStatus(ctx, objectID, userID)

	return
}

func (vr *VoteRepo) voteCancel(ctx context.Context, objectID string, userID, objectUserID string, actions []string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		for _, action := range actions {
			var (
				existsActivity entity.Activity
				has            bool
				triggerUserID,
				activityUserID string
				activityType,
				deltaRank, hasRank int
			)
			result = nil

			activityUserID, activityType, deltaRank, hasRank, err = vr.CheckRank(ctx, objectID, objectUserID, userID, action)
			if err != nil {
				return
			}

			triggerUserID = userID
			if userID == activityUserID {
				triggerUserID = "0"
			}

			has, err = session.
				Where(builder.Eq{"user_id": activityUserID}).
				And(builder.Eq{"trigger_user_id": triggerUserID}).
				And(builder.Eq{"activity_type": activityType}).
				And(builder.Eq{"object_id": objectID}).
				Get(&existsActivity)

			if !has {
				return
			}

			if existsActivity.Cancelled == entity.ActivityCancelled {
				return
			}

			if _, err = session.Where("id = ?", existsActivity.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{
					Cancelled:   entity.ActivityCancelled,
					CancelledAt: time.Now(),
				}); err != nil {
				return
			}

			// trigger user rank and send notification
			if hasRank != 0 {
				_, err = vr.userRankRepo.TriggerUserRank(ctx, session, activityUserID, -deltaRank, activityType)
				if err != nil {
					return
				}

				vr.sendNotification(ctx, activityUserID, objectUserID, objectID)
			}

			// update votes
			if action == "vote_down" || action == "vote_up" {
				votes := -1
				if action == "vote_down" {
					votes = 1
				}
				err = vr.updateVotes(ctx, session, objectID, votes)
				if err != nil {
					return
				}
			}
		}

		return
	})
	if err != nil {
		return
	}
	resp, err = vr.GetVoteResultByObjectId(ctx, objectID)
	resp.VoteStatus = vr.voteCommon.GetVoteStatus(ctx, objectID, userID)
	return
}

func (vr *VoteRepo) VoteUp(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}

	actions, ok := LimitUpActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	_, _ = vr.VoteDownCancel(ctx, objectID, userID, objectUserID)
	return vr.vote(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteDown(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitDownActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	_, _ = vr.VoteUpCancel(ctx, objectID, userID, objectUserID)
	return vr.vote(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteUpCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	var objectType string
	resp = &schema.VoteResp{}

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitUpActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	return vr.voteCancel(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteDownCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	var objectType string
	resp = &schema.VoteResp{}

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitDownActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	return vr.voteCancel(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) CheckRank(ctx context.Context, objectID, objectUserID, userID string, action string) (activityUserID string, activityType, rank, hasRank int, err error) {
	activityType, rank, hasRank, err = vr.activityRepo.GetActivityTypeByObjID(ctx, objectID, action)

	if err != nil {
		return
	}

	activityUserID = userID
	if strings.Contains(action, "voted") {
		activityUserID = objectUserID
	}

	return activityUserID, activityType, rank, hasRank, nil
}

func (vr *VoteRepo) GetVoteResultByObjectId(ctx context.Context, objectID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	for _, action := range []string{"vote_up", "vote_down"} {
		var (
			activity     entity.Activity
			votes        int64
			activityType int
		)

		activityType, _, _, _ = vr.activityRepo.GetActivityTypeByObjID(ctx, objectID, action)

		votes, err = vr.data.DB.Where(builder.Eq{"object_id": objectID}).
			And(builder.Eq{"activity_type": activityType}).
			And(builder.Eq{"cancelled": 0}).
			Count(&activity)

		if err != nil {
			return
		}

		if action == "vote_up" {
			resp.UpVotes = int(votes)
		} else {
			resp.DownVotes = int(votes)
		}
	}

	resp.Votes = resp.UpVotes - resp.DownVotes

	return resp, nil
}

func (vr *VoteRepo) ListUserVotes(
	ctx context.Context,
	userID string,
	req schema.GetVoteWithPageReq,
	activityTypes []int,
) (voteList []entity.Activity, total int64, err error) {
	session := vr.data.DB.NewSession()
	cond := builder.
		And(
			builder.Eq{"user_id": userID},
			builder.Eq{"cancelled": 0},
			builder.In("activity_type", activityTypes),
		)

	session.Where(cond).OrderBy("updated_at desc")

	total, err = pager.Help(req.Page, req.PageSize, &voteList, &entity.Activity{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// updateVotes
// if votes < 0 Decr object vote_count,otherwise Incr object vote_count
func (vr *VoteRepo) updateVotes(ctx context.Context, session *xorm.Session, objectID string, votes int) (err error) {
	var (
		objectType string
		e          error
	)

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	switch objectType {
	case "question":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Question{})
	case "answer":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Answer{})
	case "comment":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Comment{})
	default:
		e = errors.BadRequest(reason.DisallowVote)
	}

	if e != nil {
		err = e
	} else if err != nil {
		err = errors.BadRequest(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}

// sendNotification send rank triggered notification
func (vr *VoteRepo) sendNotification(ctx context.Context, activityUserID, objectUserID, objectID string) {
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		return
	}

	msg := &schema.NotificationMsg{
		ReceiverUserID: activityUserID,
		TriggerUserID:  objectUserID,
		Type:           schema.NotificationTypeAchievement,
		ObjectID:       objectID,
		ObjectType:     objectType,
	}
	notice_queue.AddNotification(msg)
}
