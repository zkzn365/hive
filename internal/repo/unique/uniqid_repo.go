package unique

import (
	"context"
	"fmt"

	"answer/internal/base/constant"
	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/service/unique"

	"github.com/segmentfault/pacman/errors"
)

// uniqueIDRepo Unique id repository
type uniqueIDRepo struct {
	data *data.Data
}

// NewUniqueIDRepo new repository
func NewUniqueIDRepo(data *data.Data) unique.UniqueIDRepo {
	return &uniqueIDRepo{
		data: data,
	}
}

// GenUniqueIDStr generate unique id string
// 1 + 00x(objectType) + 000000000000x(id)
func (ur *uniqueIDRepo) GenUniqueIDStr(ctx context.Context, key string) (uniqueID string, err error) {
	objectType := constant.ObjectTypeStrMapping[key]
	bean := &entity.Uniqid{UniqidType: objectType}
	_, err = ur.data.DB.Insert(bean)
	if err != nil {
		return "", errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return fmt.Sprintf("1%03d%013d", objectType, bean.ID), nil
}
