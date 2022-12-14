package export

import (
	"context"
	"time"

	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/service/export"

	"github.com/segmentfault/pacman/errors"
)

// emailRepo email repository
type emailRepo struct {
	data *data.Data
}

// NewEmailRepo new repository
func NewEmailRepo(data *data.Data) export.EmailRepo {
	return &emailRepo{
		data: data,
	}
}

// SetCode The email code is used to verify that the link in the message is out of date
func (e *emailRepo) SetCode(ctx context.Context, code, content string) error {
	err := e.data.Cache.SetString(ctx, code, content, 10*time.Minute)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

// VerifyCode verify the code if out of date
func (e *emailRepo) VerifyCode(ctx context.Context, code string) (content string, err error) {
	content, err = e.data.Cache.GetString(ctx, code)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
