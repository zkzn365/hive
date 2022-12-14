package activity

import (
	"context"
	"fmt"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/repo/config"
	"answer/internal/service/activity"

	"github.com/segmentfault/pacman/errors"
)

// activityRepo activity repository
type activityRepo struct {
	data *data.Data
}

// NewActivityRepo new repository
func NewActivityRepo(
	data *data.Data,
) activity.ActivityRepo {
	return &activityRepo{
		data: data,
	}
}

func (ar *activityRepo) GetObjectAllActivity(ctx context.Context, objectID string, showVote bool) (
	activityList []*entity.Activity, err error) {
	activityList = make([]*entity.Activity, 0)
	session := ar.data.DB.Desc("created_at")

	if !showVote {
		var activityTypeNotShown []int
		for _, obj := range []string{constant.AnswerObjectType, constant.QuestionObjectType, constant.CommentObjectType} {
			for _, act := range []string{
				constant.ActVotedDown,
				constant.ActVotedUp,
				constant.ActVoteDown,
				constant.ActVoteUp,
			} {
				activityTypeNotShown = append(activityTypeNotShown, config.Key2IDMapping[fmt.Sprintf("%s.%s", obj, act)])
			}
		}
		session.NotIn("activity_type", activityTypeNotShown)
	}
	err = session.Find(&activityList, &entity.Activity{OriginalObjectID: objectID})
	if err != nil {
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return activityList, nil
}
