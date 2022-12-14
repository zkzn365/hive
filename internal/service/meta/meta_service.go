package meta

import (
	"context"

	"answer/internal/base/reason"
	"answer/internal/entity"

	"github.com/segmentfault/pacman/errors"
)

// MetaRepo meta repository
type MetaRepo interface {
	AddMeta(ctx context.Context, meta *entity.Meta) (err error)
	RemoveMeta(ctx context.Context, id int) (err error)
	UpdateMeta(ctx context.Context, meta *entity.Meta) (err error)
	GetMetaByObjectIdAndKey(ctx context.Context, objectId, key string) (meta *entity.Meta, exist bool, err error)
	GetMetaList(ctx context.Context, meta *entity.Meta) (metas []*entity.Meta, err error)
}

// MetaService user service
type MetaService struct {
	metaRepo MetaRepo
}

func NewMetaService(metaRepo MetaRepo) *MetaService {
	return &MetaService{
		metaRepo: metaRepo,
	}
}

// AddMeta add meta
func (ms *MetaService) AddMeta(ctx context.Context, objID, key, value string) (err error) {
	meta := &entity.Meta{
		ObjectID: objID,
		Key:      key,
		Value:    value,
	}
	return ms.metaRepo.AddMeta(ctx, meta)
}

// RemoveMeta delete meta
func (ms *MetaService) RemoveMeta(ctx context.Context, id int) (err error) {
	return ms.metaRepo.RemoveMeta(ctx, id)
}

// UpdateMeta update meta
func (ms *MetaService) UpdateMeta(ctx context.Context, metaID int, key, value string) (err error) {
	meta := &entity.Meta{
		ID:    metaID,
		Key:   key,
		Value: value,
	}
	return ms.metaRepo.UpdateMeta(ctx, meta)
}

// GetMetaByObjectIdAndKey get meta one
func (ms *MetaService) GetMetaByObjectIdAndKey(ctx context.Context, objectID, key string) (meta *entity.Meta, err error) {
	meta, exist, err := ms.metaRepo.GetMetaByObjectIdAndKey(ctx, objectID, key)
	if err != nil {
		return
	}
	if !exist {
		return nil, errors.BadRequest(reason.UnknownError)
	}
	return meta, nil
}

// GetMetaList get meta list all
func (ms *MetaService) GetMetaList(ctx context.Context, objID string) (metas []*entity.Meta, err error) {
	metas, err = ms.metaRepo.GetMetaList(ctx, &entity.Meta{ObjectID: objID})
	if err != nil {
		return nil, err
	}
	return metas, err
}
