package report

import (
	"encoding/json"

	"answer/internal/base/constant"
	"answer/internal/base/reason"
	"answer/internal/base/translator"
	"answer/internal/entity"
	"answer/internal/schema"
	"answer/internal/service/object_info"
	"answer/internal/service/report_common"
	"answer/pkg/obj"

	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/i18n"
	"golang.org/x/net/context"
)

// ReportService user service
type ReportService struct {
	reportRepo        report_common.ReportRepo
	objectInfoService *object_info.ObjService
}

// NewReportService new report service
func NewReportService(reportRepo report_common.ReportRepo,
	objectInfoService *object_info.ObjService,
) *ReportService {
	return &ReportService{
		reportRepo:        reportRepo,
		objectInfoService: objectInfoService,
	}
}

// AddReport add report
func (rs *ReportService) AddReport(ctx context.Context, req *schema.AddReportReq) (err error) {
	objectTypeNumber, err := obj.GetObjectTypeNumberByObjectID(req.ObjectID)
	if err != nil {
		return err
	}

	// TODO this reported user id should be get by revision
	objInfo, err := rs.objectInfoService.GetInfo(ctx, req.ObjectID)
	if err != nil {
		return err
	}

	report := &entity.Report{
		UserID:         req.UserID,
		ReportedUserID: objInfo.ObjectCreatorUserID,
		ObjectID:       req.ObjectID,
		ObjectType:     objectTypeNumber,
		ReportType:     req.ReportType,
		Content:        req.Content,
		Status:         entity.ReportStatusPending,
	}
	return rs.reportRepo.AddReport(ctx, report)
}

// GetReportTypeList get report list all
func (rs *ReportService) GetReportTypeList(ctx context.Context, lang i18n.Language, req *schema.GetReportListReq) (
	resp []*schema.GetReportTypeResp, err error,
) {
	resp = make([]*schema.GetReportTypeResp, 0)
	switch req.Source {
	case constant.QuestionObjectType:
		err = json.Unmarshal([]byte(constant.QuestionReportJSON), &resp)
	case constant.AnswerObjectType:
		err = json.Unmarshal([]byte(constant.AnswerReportJSON), &resp)
	case constant.CommentObjectType:
		err = json.Unmarshal([]byte(constant.CommentReportJSON), &resp)
	}
	if err != nil {
		err = errors.BadRequest(reason.UnknownError)
	}
	for _, t := range resp {
		t.Name = translator.GlobalTrans.Tr(lang, t.Name)
		t.Description = translator.GlobalTrans.Tr(lang, t.Description)
	}
	return resp, err
}
