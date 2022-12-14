package report

import (
	"context"

	"answer/internal/base/constant"
	"answer/internal/base/pager"
	"answer/internal/schema"
	"answer/internal/service/report_common"

	"answer/internal/base/data"
	"answer/internal/base/reason"
	"answer/internal/entity"
	"answer/internal/service/unique"

	"github.com/segmentfault/pacman/errors"
)

// reportRepo report repository
type reportRepo struct {
	data         *data.Data
	uniqueIDRepo unique.UniqueIDRepo
}

// NewReportRepo new repository
func NewReportRepo(data *data.Data, uniqueIDRepo unique.UniqueIDRepo) report_common.ReportRepo {
	return &reportRepo{
		data:         data,
		uniqueIDRepo: uniqueIDRepo,
	}
}

// AddReport add report
func (rr *reportRepo) AddReport(ctx context.Context, report *entity.Report) (err error) {
	report.ID, err = rr.uniqueIDRepo.GenUniqueIDStr(ctx, report.TableName())
	if err != nil {
		return err
	}
	_, err = rr.data.DB.Insert(report)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetReportListPage get report list page
func (rr *reportRepo) GetReportListPage(ctx context.Context, dto schema.GetReportListPageDTO) (reports []entity.Report, total int64, err error) {
	var (
		ok         bool
		status     int
		objectType int
		session    = rr.data.DB.NewSession()
		cond       = entity.Report{}
	)

	// parse status
	status, ok = entity.ReportStatus[dto.Status]
	if !ok {
		status = entity.ReportStatus["pending"]
	}
	cond.Status = status

	// parse object type
	objectType, ok = constant.ObjectTypeStrMapping[dto.ObjectType]
	if ok {
		cond.ObjectType = objectType
	}

	// order
	session.OrderBy("updated_at desc")

	total, err = pager.Help(dto.Page, dto.PageSize, &reports, cond, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

// GetByID get report by ID
func (ar *reportRepo) GetByID(ctx context.Context, id string) (report entity.Report, exist bool, err error) {
	report = entity.Report{}
	exist, err = ar.data.DB.ID(id).Get(&report)
	return
}

// UpdateByID handle report by ID
func (ar *reportRepo) UpdateByID(
	ctx context.Context,
	id string,
	handleData entity.Report,
) (err error) {
	_, err = ar.data.DB.ID(id).Update(&handleData)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (vr *reportRepo) GetReportCount(ctx context.Context) (count int64, err error) {
	list := make([]*entity.Report, 0)
	count, err = vr.data.DB.Where("status =?", entity.ReportStatusPending).FindAndCount(&list)
	if err != nil {
		return count, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}
