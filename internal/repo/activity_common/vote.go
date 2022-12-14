package activity_common

import (
	"context"

	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/service/activity_common"

	"github.com/segmentfault/pacman/errors"
)

// VoteRepo activity repository
type VoteRepo struct {
	data         *data.Data
	activityRepo activity_common.ActivityRepo
}

// NewVoteRepo new repository
func NewVoteRepo(data *data.Data, activityRepo activity_common.ActivityRepo) activity_common.VoteRepo {
	return &VoteRepo{
		data:         data,
		activityRepo: activityRepo,
	}
}

func (vr *VoteRepo) GetVoteStatus(ctx context.Context, objectID, userID string) (status string) {
	for _, action := range []string{"vote_up", "vote_down"} {
		at := &entity.Activity{}
		activityType, _, _, err := vr.activityRepo.GetActivityTypeByObjID(ctx, objectID, action)
		if err != nil {
			return ""
		}
		has, err := vr.data.DB.Where("object_id =? AND cancelled=0 AND activity_type=? AND user_id=?", objectID, activityType, userID).Get(at)
		if err != nil {
			return ""
		}
		if has {
			return action
		}
	}
	return ""
}

func (vr *VoteRepo) GetVoteCount(ctx context.Context, activityTypes []int) (count int64, err error) {
	list := make([]*entity.Activity, 0)
	count, err = vr.data.DB.Where("cancelled =0").In("activity_type", activityTypes).FindAndCount(&list)
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
