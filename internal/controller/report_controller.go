package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/base/reason"
	"answer/internal/schema"
	"answer/internal/service/rank"
	"answer/internal/service/report"

	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

// ReportController report controller
type ReportController struct {
	reportService *report.ReportService
	rankService   *rank.RankService
}

// NewReportController new controller
func NewReportController(reportService *report.ReportService, rankService *rank.RankService) *ReportController {
	return &ReportController{reportService: reportService, rankService: rankService}
}

// AddReport add report
// @Summary add report
// @Description add report <br> source (question, answer, comment, user)
// @Security ApiKeyAuth
// @Tags Report
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.AddReportReq true "report"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/report [post]
func (rc *ReportController) AddReport(ctx *gin.Context) {
	req := &schema.AddReportReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := rc.rankService.CheckOperationPermission(ctx, req.UserID, rank.ReportAddRank, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = rc.reportService.AddReport(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetReportTypeList get report type list
// @Summary get report type list
// @Description get report type list
// @Tags Report
// @Produce json
// @Param source query string true "report source" Enums(question, answer, comment, user)
// @Success 200 {object} handler.RespBody{data=[]schema.GetReportTypeResp}
// @Router /answer/api/v1/report/type/list [get]
func (rc *ReportController) GetReportTypeList(ctx *gin.Context) {
	req := &schema.GetReportListReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	resp, err := rc.reportService.GetReportTypeList(ctx, handler.GetLang(ctx), req)
	handler.HandleResponse(ctx, err, resp)
}
