package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/schema"
	"answer/internal/service/notification"
	"answer/internal/service/rank"

	"github.com/gin-gonic/gin"
)

// NotificationController notification controller
type NotificationController struct {
	notificationService *notification.NotificationService
	rankService         *rank.RankService
}

// NewNotificationController new controller
func NewNotificationController(
	notificationService *notification.NotificationService,
	rankService *rank.RankService,
) *NotificationController {
	return &NotificationController{
		notificationService: notificationService,
		rankService:         rankService,
	}
}

// GetRedDot
// @Summary GetRedDot
// @Description GetRedDot
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/notification/status [get]
func (nc *NotificationController) GetRedDot(ctx *gin.Context) {

	req := &schema.GetRedDot{}

	userID := middleware.GetLoginUserIDFromContext(ctx)
	req.UserID = userID
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := nc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		rank.QuestionAuditRank,
		rank.AnswerAuditRank,
		rank.TagAuditRank,
	}, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanReviewQuestion = canList[0]
	req.CanReviewAnswer = canList[1]
	req.CanReviewTag = canList[2]

	RedDot, err := nc.notificationService.GetRedDot(ctx, req)
	handler.HandleResponse(ctx, err, RedDot)
}

// ClearRedDot
// @Summary DelRedDot
// @Description DelRedDot
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.NotificationClearRequest true "NotificationClearRequest"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/notification/status [put]
func (nc *NotificationController) ClearRedDot(ctx *gin.Context) {
	req := &schema.NotificationClearRequest{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := nc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		rank.QuestionAuditRank,
		rank.AnswerAuditRank,
		rank.TagAuditRank,
	}, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanReviewQuestion = canList[0]
	req.CanReviewAnswer = canList[1]
	req.CanReviewTag = canList[2]

	RedDot, err := nc.notificationService.ClearRedDot(ctx, req)
	handler.HandleResponse(ctx, err, RedDot)
}

// ClearUnRead
// @Summary ClearUnRead
// @Description ClearUnRead
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.NotificationClearRequest true "NotificationClearRequest"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/notification/read/state/all [put]
func (nc *NotificationController) ClearUnRead(ctx *gin.Context) {
	req := &schema.NotificationClearRequest{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	userID := middleware.GetLoginUserIDFromContext(ctx)
	err := nc.notificationService.ClearUnRead(ctx, userID, req.TypeStr)
	handler.HandleResponse(ctx, err, gin.H{})
}

// ClearIDUnRead
// @Summary ClearUnRead
// @Description ClearUnRead
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.NotificationClearIDRequest true "NotificationClearIDRequest"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/notification/read/state [put]
func (nc *NotificationController) ClearIDUnRead(ctx *gin.Context) {
	req := &schema.NotificationClearIDRequest{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	userID := middleware.GetLoginUserIDFromContext(ctx)
	err := nc.notificationService.ClearIDUnRead(ctx, userID, req.ID)
	handler.HandleResponse(ctx, err, gin.H{})
}

// GetList get notification list
// @Summary get notification list
// @Description get notification list
// @Tags Notification
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param type query string true "type" Enums(inbox,achievement)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/notification/page [get]
func (nc *NotificationController) GetList(ctx *gin.Context) {
	req := &schema.NotificationSearch{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	resp, err := nc.notificationService.GetNotificationPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
