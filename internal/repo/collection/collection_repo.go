package collection

import (
	"context"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/pager"
	"answer/internal/base/reason"
	"answer/internal/entity"
	collectioncommon "answer/internal/service/collection_common"
	"answer/internal/service/unique"

	"github.com/segmentfault/pacman/errors"
)

// collectionRepo collection repository
type collectionRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewCollectionRepo new repository
func NewCollectionRepo(data *data.Data, uniqueIDRepo unique.UniqueIDRepo) collectioncommon.CollectionRepo {
	return &collectionRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddCollection add collection
func (cr *collectionRepo) AddCollection(ctx context.Context, collection *entity.Collection) (err error) {
	id, err := cr.uniqueIDRepo.GenUniqueIDStr(ctx, collection.TableName())
	if err == nil {
		collection.ID = id
		_, err = cr.data.DB.Insert(collection)
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}
	return nil
}

// RemoveCollection delete collection
func (cr *collectionRepo) RemoveCollection(ctx context.Context, id string) (err error) {
	_, err = cr.data.DB.Where("id =?", id).Delete(&entity.Collection{})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// UpdateCollection update collection
func (cr *collectionRepo) UpdateCollection(ctx context.Context, collection *entity.Collection, cols []string) (err error) {
	_, err = cr.data.DB.ID(collection.ID).Cols(cols...).Update(collection)
	return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
}

// GetCollection get collection one
func (cr *collectionRepo) GetCollection(ctx context.Context, id int) (collection *entity.Collection, exist bool, err error) {
	collection = &entity.Collection{}
	exist, err = cr.data.DB.ID(id).Get(collection)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetCollectionList get collection list all
func (cr *collectionRepo) GetCollectionList(ctx context.Context, collection *entity.Collection) (collectionList []*entity.Collection, err error) {
	collectionList = make([]*entity.Collection, 0)
	err = cr.data.DB.Find(collectionList, collection)
	err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	return
}

// GetOneByObjectIDAndUser get one by object TagID and user
func (cr *collectionRepo) GetOneByObjectIDAndUser(ctx context.Context, userID string, objectID string) (collection *entity.Collection, exist bool, err error) {
	collection = &entity.Collection{}
	exist, err = cr.data.DB.Where("user_id = ? and object_id = ?", userID, objectID).Get(collection)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// SearchByObjectIDsAndUser search by object IDs and user
func (cr *collectionRepo) SearchByObjectIDsAndUser(ctx context.Context, userID string, objectIDs []string) ([]*entity.Collection, error) {
	collectionList := make([]*entity.Collection, 0)
	err := cr.data.DB.Where("user_id = ?", userID).In("object_id", objectIDs).Find(&collectionList)
	if err != nil {
		return collectionList, err
	}
	return collectionList, nil
}

// CountByObjectID count by object TagID
func (cr *collectionRepo) CountByObjectID(ctx context.Context, objectID string) (total int64, err error) {
	collection := &entity.Collection{}
	total, err = cr.data.DB.Where("object_id = ?", objectID).Count(collection)
	if err != nil {
		return 0, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetCollectionPage get collection page
func (cr *collectionRepo) GetCollectionPage(ctx context.Context, page, pageSize int, collection *entity.Collection) (collectionList []*entity.Collection, total int64, err error) {
	collectionList = make([]*entity.Collection, 0)

	session := cr.data.DB.NewSession()
	if collection.UserID != "" && collection.UserID != "0" {
		session = session.Where("user_id = ?", collection.UserID)
	}

	if collection.UserCollectionGroupID != "" && collection.UserCollectionGroupID != "0" {
		session = session.Where("user_collection_group_id = ?", collection.UserCollectionGroupID)
	}
	session = session.OrderBy("update_time desc")

	total, err = pager.Help(page, pageSize, collectionList, collection, session)
	err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	return
}

// SearchObjectCollected check object is collected or not
func (cr *collectionRepo) SearchObjectCollected(ctx context.Context, userID string, objectIds []string) (map[string]bool, error) {
	collectedMap := make(map[string]bool)
	list, err := cr.SearchByObjectIDsAndUser(ctx, userID, objectIds)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		return collectedMap, err
	}
	for _, item := range list {
		collectedMap[item.ObjectID] = true
	}
	return collectedMap, err
}

// SearchList
func (cr *collectionRepo) SearchList(ctx context.Context, search *entity.CollectionSearch) ([]*entity.Collection, int64, error) {
	var count int64
	var err error
	rows := make([]*entity.Collection, 0)
	if search.Page > 0 {
		search.Page = search.Page - 1
	} else {
		search.Page = 0
	}
	if search.PageSize == 0 {
		search.PageSize = constant.DefaultPageSize
	}
	offset := search.Page * search.PageSize
	session := cr.data.DB.Where("")
	if len(search.UserID) > 0 {
		session = session.And("user_id = ?", search.UserID)
	} else {
		return rows, count, nil
	}
	session = session.Limit(search.PageSize, offset)
	count, err = session.OrderBy("updated_at desc").FindAndCount(&rows)
	if err != nil {
		return rows, count, err
	}
	return rows, count, nil
}
