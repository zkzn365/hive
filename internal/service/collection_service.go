package service

import (
	"context"
	"fmt"

	"answer/internal/entity"
	"answer/internal/schema"
	collectioncommon "answer/internal/service/collection_common"
	questioncommon "answer/internal/service/question_common"

	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
)

// CollectionService user service
type CollectionService struct {
	collectionRepo      collectioncommon.CollectionRepo
	collectionGroupRepo CollectionGroupRepo
	questionCommon      *questioncommon.QuestionCommon
}

func NewCollectionService(
	collectionRepo collectioncommon.CollectionRepo,
	collectionGroupRepo CollectionGroupRepo,
	questionCommon *questioncommon.QuestionCommon,
) *CollectionService {
	return &CollectionService{
		collectionRepo:      collectionRepo,
		collectionGroupRepo: collectionGroupRepo,
		questionCommon:      questionCommon,
	}
}

func (cs *CollectionService) CollectionSwitch(ctx context.Context, dto *schema.CollectionSwitchDTO) (resp *schema.CollectionSwitchResp, err error) {
	resp = &schema.CollectionSwitchResp{}
	dbData, has, err := cs.collectionRepo.GetOneByObjectIDAndUser(ctx, dto.UserID, dto.ObjectID)
	if err != nil {
		return
	}
	if has {
		err = cs.collectionRepo.RemoveCollection(ctx, dbData.ID)
		if err != nil {
			return nil, err
		}
		err = cs.questionCommon.UpdateCollectionCount(ctx, dto.ObjectID, -1)
		if err != nil {
			log.Error("UpdateCollectionCount", err.Error())
		}
		var count int64
		count, err = cs.objectCollectionCount(ctx, dto.ObjectID)
		if err != nil {
			return resp, err
		}
		resp.ObjectCollectionCount = fmt.Sprintf("%v", count)
		resp.Switch = false
		return resp, err
	}

	if dto.GroupID == "" || dto.GroupID == "0" {
		var (
			defaultGroup *entity.CollectionGroup
			has          bool
		)
		defaultGroup, has, err = cs.collectionGroupRepo.GetDefaultID(ctx, dto.UserID)
		if err != nil {
			return nil, err
		}
		if !has {
			var dbdefaultGroup *entity.CollectionGroup
			dbdefaultGroup, err = cs.collectionGroupRepo.AddCollectionDefaultGroup(ctx, dto.UserID)
			if err != nil {
				return nil, err
			}
			dto.GroupID = dbdefaultGroup.ID
		} else {
			dto.GroupID = defaultGroup.ID
		}
	}
	collection := &entity.Collection{
		UserCollectionGroupID: dto.GroupID,
		UserID:                dto.UserID,
		ObjectID:              dto.ObjectID,
	}

	err = cs.collectionRepo.AddCollection(ctx, collection)
	if err != nil {
		return
	}
	err = cs.questionCommon.UpdateCollectionCount(ctx, dto.ObjectID, 1)
	if err != nil {
		log.Error("UpdateCollectionCount", err.Error())
	}
	count, err := cs.objectCollectionCount(ctx, dto.ObjectID)
	if err != nil {
		return
	}
	resp.ObjectCollectionCount = fmt.Sprintf("%d", count)
	resp.Switch = true
	return
}

func (cs *CollectionService) objectCollectionCount(ctx context.Context, objectID string) (int64, error) {
	count, err := cs.collectionRepo.CountByObjectID(ctx, objectID)
	return count, err
}

func (cs *CollectionService) add(ctx context.Context, collection *entity.Collection) error {
	_, has, err := cs.collectionRepo.GetOneByObjectIDAndUser(ctx, collection.UserID, collection.ObjectID)
	if err != nil {
		return err
	}
	if has {
		return errors.BadRequest("already collected")
	}

	if collection.UserCollectionGroupID == "" || collection.UserCollectionGroupID == "0" {
		var (
			defaultGroup *entity.CollectionGroup
			has          bool
		)
		defaultGroup, has, err = cs.collectionGroupRepo.GetDefaultID(ctx, collection.UserID)
		if err != nil {
			return err
		}
		if !has {
			defaultGroup, err = cs.collectionGroupRepo.AddCollectionDefaultGroup(ctx, collection.UserID)
			if err != nil {
				return err
			}
			collection.UserCollectionGroupID = defaultGroup.ID

		} else {
			collection.UserCollectionGroupID = defaultGroup.ID
		}
	}
	err = cs.collectionRepo.AddCollection(ctx, collection)
	if err != nil {
		return err
	}
	return nil
}

// Cancel
func (cs *CollectionService) cancel(ctx context.Context, collection *entity.Collection) error {
	dbData, has, err := cs.collectionRepo.GetOneByObjectIDAndUser(ctx, collection.UserID, collection.ObjectID)
	if err != nil {
		return err
	}
	if !has {
		return errors.BadRequest("collected record does not exist")
	}
	err = cs.collectionRepo.RemoveCollection(ctx, dbData.ID)
	if err != nil {
		return err
	}
	return nil
}
