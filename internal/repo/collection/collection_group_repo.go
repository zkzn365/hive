package collection

import (
	"context"

	"answer/internal/base/data"
	"answer/internal/base/pager"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service"

	"github.com/segmentfault/pacman/errors"
)

// collectionGroupRepo collectionGroup repository
type collectionGroupRepo struct {
	data *data.Data
}

// NewCollectionGroupRepo new repository
func NewCollectionGroupRepo(data *data.Data) service.CollectionGroupRepo {
	return &collectionGroupRepo{
		data: data,
	}
}

// AddCollectionGroup add collection group
func (cr *collectionGroupRepo) AddCollectionGroup(ctx context.Context, collectionGroup *entity.CollectionGroup) (err error) {
	_, err = cr.data.DB.Insert(collectionGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// AddCollectionDefaultGroup add collection group
func (cr *collectionGroupRepo) AddCollectionDefaultGroup(ctx context.Context, userID string) (collectionGroup *entity.CollectionGroup, err error) {
	defaultGroup := &entity.CollectionGroup{
		Name:         "default",
		DefaultGroup: schema.CGDefault,
		UserID:       userID,
	}
	_, err = cr.data.DB.Insert(defaultGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return
	}
	collectionGroup = defaultGroup
	return
}

// UpdateCollectionGroup update collection group
func (cr *collectionGroupRepo) UpdateCollectionGroup(ctx context.Context, collectionGroup *entity.CollectionGroup, cols []string) (err error) {
	_, err = cr.data.DB.ID(collectionGroup.ID).Cols(cols...).Update(collectionGroup)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetCollectionGroup get collection group one
func (cr *collectionGroupRepo) GetCollectionGroup(ctx context.Context, id string) (
	collectionGroup *entity.CollectionGroup, exist bool, err error,
) {
	collectionGroup = &entity.CollectionGroup{}
	exist, err = cr.data.DB.ID(id).Get(collectionGroup)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetCollectionGroupPage get collection group page
func (cr *collectionGroupRepo) GetCollectionGroupPage(ctx context.Context, page, pageSize int, collectionGroup *entity.CollectionGroup) (collectionGroupList []*entity.CollectionGroup, total int64, err error) {
	collectionGroupList = make([]*entity.CollectionGroup, 0)

	session := cr.data.DB.NewSession()
	if collectionGroup.UserID != "" && collectionGroup.UserID != "0" {
		session = session.Where("user_id = ?", collectionGroup.UserID)
	}
	session = session.OrderBy("update_time desc")

	total, err = pager.Help(page, pageSize, collectionGroupList, collectionGroup, session)
	err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	return
}

func (cr *collectionGroupRepo) GetDefaultID(ctx context.Context, userID string) (collectionGroup *entity.CollectionGroup, has bool, err error) {
	collectionGroup = &entity.CollectionGroup{}
	has, err = cr.data.DB.Where("user_id =? and  default_group = ?", userID, schema.CGDefault).Get(collectionGroup)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return
	}
	return
}
