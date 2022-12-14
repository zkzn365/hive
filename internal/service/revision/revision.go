package revision

import (
	"context"

	"answer/internal/entity"

	"xorm.io/xorm"
)

// RevisionRepo revision repository
type RevisionRepo interface {
	AddRevision(ctx context.Context, revision *entity.Revision, autoUpdateRevisionID bool) (err error)
	GetRevisionByID(ctx context.Context, revisionID string) (revision *entity.Revision, exist bool, err error)
	GetLastRevisionByObjectID(ctx context.Context, objectID string) (revision *entity.Revision, exist bool, err error)
	GetRevisionList(ctx context.Context, revision *entity.Revision) (revisionList []entity.Revision, err error)
	UpdateObjectRevisionId(ctx context.Context, revision *entity.Revision, session *xorm.Session) (err error)
	ExistUnreviewedByObjectID(ctx context.Context, objectID string) (revision *entity.Revision, exist bool, err error)
	GetUnreviewedRevisionPage(ctx context.Context, page, pageSize int, objectTypes []int) ([]*entity.Revision, int64, error)
	UpdateStatus(ctx context.Context, id string, status int, reviewUserID string) (err error)
}
