package report_common

import (
	"context"

	"answer/internal/entity"
	"answer/internal/schema"
)

// ReportRepo report repository
type ReportRepo interface {
	AddReport(ctx context.Context, report *entity.Report) (err error)
	GetReportListPage(ctx context.Context, query schema.GetReportListPageDTO) (reports []entity.Report, total int64, err error)
	GetByID(ctx context.Context, id string) (report entity.Report, exist bool, err error)
	UpdateByID(ctx context.Context, id string, handleData entity.Report) (err error)
	GetReportCount(ctx context.Context) (count int64, err error)
}
