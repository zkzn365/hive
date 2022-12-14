package rank

import (
	"context"

	"answer/internal/base/data"
	"answer/internal/base/pager"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/service/config"
	"answer/internal/service/rank"

	"github.com/jinzhu/now"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// UserRankRepo user rank repository
type UserRankRepo struct {
	data       *data.Data
	configRepo config.ConfigRepo
}

// NewUserRankRepo new repository
func NewUserRankRepo(data *data.Data, configRepo config.ConfigRepo) rank.UserRankRepo {
	return &UserRankRepo{
		data:       data,
		configRepo: configRepo,
	}
}

// TriggerUserRank trigger user rank change
// session is need provider, it means this action must be success or failure
// if outer action is failed then this action is need rollback
func (ur *UserRankRepo) TriggerUserRank(ctx context.Context,
	session *xorm.Session, userID string, deltaRank int, activityType int,
) (isReachStandard bool, err error) {
	if deltaRank == 0 {
		return false, nil
	}

	if deltaRank < 0 {
		// if user rank is lower than 1 after this action, then user rank will be set to 1 only.
		var isReachMin bool
		isReachMin, err = ur.checkUserMinRank(ctx, session, userID, activityType)
		if err != nil {
			return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if isReachMin {
			_, err = session.Where(builder.Eq{"id": userID}).Update(&entity.User{Rank: 1})
			if err != nil {
				return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			return false, nil
		}
	} else {
		isReachStandard, err = ur.checkUserTodayRank(ctx, session, userID, activityType)
		if err != nil {
			return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if isReachStandard {
			return isReachStandard, nil
		}
	}
	_, err = session.Where(builder.Eq{"id": userID}).Incr("`rank`", deltaRank).Update(&entity.User{})
	if err != nil {
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return false, nil
}

func (ur *UserRankRepo) checkUserMinRank(ctx context.Context, session *xorm.Session, userID string, deltaRank int) (
	isReachStandard bool, err error,
) {
	bean := &entity.User{ID: userID}
	_, err = session.Select("rank").Get(bean)
	if err != nil {
		return false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if bean.Rank+deltaRank < 1 {
		log.Infof("user %s is rank %d out of range before rank operation", userID, deltaRank)
		return true, nil
	}
	return
}

func (ur *UserRankRepo) checkUserTodayRank(ctx context.Context,
	session *xorm.Session, userID string, activityType int,
) (isReachStandard bool, err error) {
	// exclude daily rank
	exclude, _ := ur.configRepo.GetArrayString("daily_rank_limit.exclude")
	for _, item := range exclude {
		var excludeActivityType int
		excludeActivityType, err = ur.configRepo.GetInt(item)
		if err != nil {
			return false, err
		}
		if activityType == excludeActivityType {
			return false, nil
		}
	}

	// get user
	start, end := now.BeginningOfDay(), now.EndOfDay()
	session.Where(builder.Eq{"user_id": userID})
	session.Where(builder.Eq{"cancelled": 0})
	session.Where(builder.Between{
		Col:     "updated_at",
		LessVal: start,
		MoreVal: end,
	})
	earned, err := session.Sum(&entity.Activity{}, "rank")
	if err != nil {
		return false, err
	}

	// max rank
	maxDailyRank, err := ur.configRepo.GetInt("daily_rank_limit")
	if err != nil {
		return false, err
	}

	if int(earned) < maxDailyRank {
		return false, nil
	}
	log.Infof("user %s today has rank %d is reach stand %d", userID, earned, maxDailyRank)
	return true, nil
}

func (ur *UserRankRepo) UserRankPage(ctx context.Context, userID string, page, pageSize int) (
	rankPage []*entity.Activity, total int64, err error,
) {
	rankPage = make([]*entity.Activity, 0)

	session := ur.data.DB.Where(builder.Eq{"has_rank": 1}.And(builder.Eq{"cancelled": 0}))
	session.Desc("created_at")

	cond := &entity.Activity{UserID: userID}
	total, err = pager.Help(page, pageSize, &rankPage, cond, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
