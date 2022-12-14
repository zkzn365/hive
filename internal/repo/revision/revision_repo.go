package revision

import (
	"context"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/pager"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/service/revision"
	"answer/internal/service/unique"
	"answer/pkg/converter"
	"answer/pkg/obj"

	"github.com/segmentfault/pacman/errors"
	"xorm.io/builder"
	"xorm.io/xorm"
)

// revisionRepo revision repository
type revisionRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewRevisionRepo new repository
func NewRevisionRepo(data *data.Data, uniqueIDRepo unique.UniqueIDRepo) revision.RevisionRepo {
	return &revisionRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddRevision add revision
// autoUpdateRevisionID bool : if autoUpdateRevisionID is true , the object.revision_id will be updated,
// if not need auto update object.revision_id, it must be false.
// example: user can edit the object, but need audit, the revision_id will be updated when admin approved
func (rr *revisionRepo) AddRevision(ctx context.Context, revision *entity.Revision, autoUpdateRevisionID bool) (err error) {
	objectTypeNumber, err := obj.GetObjectTypeNumberByObjectID(revision.ObjectID)
	if err != nil {
		return errors.BadRequest(reason.ObjectNotFound)
	}

	revision.ObjectType = objectTypeNumber
	if !rr.allowRecord(revision.ObjectType) {
		return nil
	}
	_, err = rr.data.DB.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err = session.Insert(revision)
		if err != nil {
			_ = session.Rollback()
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if autoUpdateRevisionID {
			err = rr.UpdateObjectRevisionId(ctx, revision, session)
			if err != nil {
				_ = session.Rollback()
				return nil, err
			}
		}
		return nil, nil
	})

	return err
}

// UpdateObjectRevisionId updates the object.revision_id field
func (rr *revisionRepo) UpdateObjectRevisionId(ctx context.Context, revision *entity.Revision, session *xorm.Session) (err error) {
	tableName, err := obj.GetObjectTypeStrByObjectID(revision.ObjectID)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	_, err = session.Table(tableName).Where("id = ?", revision.ObjectID).Cols("`revision_id`").Update(struct {
		RevisionID string `xorm:"revision_id"`
	}{
		RevisionID: revision.ID,
	})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// UpdateStatus update revision status
func (rr *revisionRepo) UpdateStatus(ctx context.Context, id string, status int, reviewUserID string) (err error) {
	if id == "" {
		return nil
	}
	var data entity.Revision
	data.ID = id
	data.Status = status
	data.ReviewUserID = converter.StringToInt64(reviewUserID)
	_, err = rr.data.DB.Where("id =?", id).Cols("status", "review_user_id").Update(&data)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// GetRevision get revision one
func (rr *revisionRepo) GetRevision(ctx context.Context, id string) (
	revision *entity.Revision, exist bool, err error,
) {
	revision = &entity.Revision{}
	exist, err = rr.data.DB.ID(id).Get(revision)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetRevisionByID get object's last revision by object TagID
func (rr *revisionRepo) GetRevisionByID(ctx context.Context, revisionID string) (
	revision *entity.Revision, exist bool, err error) {
	revision = &entity.Revision{}
	exist, err = rr.data.DB.Where("id = ?", revisionID).Get(revision)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (rr *revisionRepo) ExistUnreviewedByObjectID(ctx context.Context, objectID string) (
	revision *entity.Revision, exist bool, err error) {
	revision = &entity.Revision{}
	exist, err = rr.data.DB.Where("object_id = ?", objectID).And("status = ?", entity.RevisionUnreviewedStatus).Get(revision)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetLastRevisionByObjectID get object's last revision by object TagID
func (rr *revisionRepo) GetLastRevisionByObjectID(ctx context.Context, objectID string) (
	revision *entity.Revision, exist bool, err error,
) {
	revision = &entity.Revision{}
	exist, err = rr.data.DB.Where("object_id = ?", objectID).OrderBy("created_at DESC").Get(revision)
	if err != nil {
		return nil, false, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetRevisionList get revision list all
func (rr *revisionRepo) GetRevisionList(ctx context.Context, revision *entity.Revision) (revisionList []entity.Revision, err error) {
	revisionList = []entity.Revision{}
	err = rr.data.DB.Where(builder.Eq{
		"object_id": revision.ObjectID,
	}).OrderBy("created_at DESC").Find(&revisionList)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// allowRecord check the object type can record revision or not
func (rr *revisionRepo) allowRecord(objectType int) (ok bool) {
	switch objectType {
	case constant.ObjectTypeStrMapping["question"]:
		return true
	case constant.ObjectTypeStrMapping["answer"]:
		return true
	case constant.ObjectTypeStrMapping["tag"]:
		return true
	default:
		return false
	}
}

// GetUnreviewedRevisionPage get unreviewed revision page
func (rr *revisionRepo) GetUnreviewedRevisionPage(ctx context.Context, page int, pageSize int,
	objectTypeList []int) (revisionList []*entity.Revision, total int64, err error) {
	revisionList = make([]*entity.Revision, 0)
	if len(objectTypeList) == 0 {
		return revisionList, 0, nil
	}
	session := rr.data.DB.NewSession()
	session = session.And("status = ?", entity.RevisionUnreviewedStatus)
	session = session.In("object_type", objectTypeList)
	session = session.OrderBy("created_at asc")

	total, err = pager.Help(page, pageSize, &revisionList, &entity.Revision{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
