package controller

import (
	"answer/internal/base/constant"
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/base/reason"
	"answer/internal/schema"
	"answer/internal/service"
	"answer/internal/service/rank"
	"answer/pkg/obj"

	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

// RevisionController revision controller
type RevisionController struct {
	revisionListService *service.RevisionService
	rankService         *rank.RankService
}

// NewRevisionController new controller
func NewRevisionController(
	revisionListService *service.RevisionService,
	rankService *rank.RankService,
) *RevisionController {
	return &RevisionController{
		revisionListService: revisionListService,
		rankService:         rankService,
	}
}

// GetRevisionList godoc
// @Summary get revision list
// @Description get revision list
// @Tags Revision
// @Produce json
// @Param object_id query string true "object id"
// @Success 200 {object} handler.RespBody{data=[]schema.GetRevisionResp}
// @Router /answer/api/v1/revisions [get]
func (rc *RevisionController) GetRevisionList(ctx *gin.Context) {
	objectID := ctx.Query("object_id")
	if objectID == "0" || objectID == "" {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}

	req := &schema.GetRevisionListReq{
		ObjectID: objectID,
	}

	resp, err := rc.revisionListService.GetRevisionList(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetUnreviewedRevisionList godoc
// @Summary get unreviewed revision list
// @Description get unreviewed revision list
// @Tags Revision
// @Produce json
// @Security ApiKeyAuth
// @Param page query string true "page id"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetUnreviewedRevisionResp}}
// @Router /answer/api/v1/revisions/unreviewed [get]
func (rc *RevisionController) GetUnreviewedRevisionList(ctx *gin.Context) {
	req := &schema.RevisionSearch{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := rc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
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

	resp, err := rc.revisionListService.GetUnreviewedRevisionPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// RevisionAudit godoc
// @Summary revision audit
// @Description revision audit operation:approve or reject
// @Tags Revision
// @Produce json
// @Security ApiKeyAuth
// @Param data body schema.RevisionAuditReq true "audit"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/api/v1/revisions/audit [put]
func (rc *RevisionController) RevisionAudit(ctx *gin.Context) {
	req := &schema.RevisionAuditReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := rc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
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

	err = rc.revisionListService.RevisionAudit(ctx, req)
	handler.HandleResponse(ctx, err, gin.H{})
}

// CheckCanUpdateRevision check can update revision
// @Summary check can update revision
// @Description check can update revision
// @Tags Revision
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query string true "id" default(string)
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/revisions/edit/check [get]
func (rc *RevisionController) CheckCanUpdateRevision(ctx *gin.Context) {
	req := &schema.CheckCanQuestionUpdate{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	action := ""
	objectTypeStr, _ := obj.GetObjectTypeStrByObjectID(req.ID)
	switch objectTypeStr {
	case constant.QuestionObjectType:
		action = rank.QuestionEditRank
	case constant.AnswerObjectType:
		action = rank.AnswerEditRank
	case constant.TagObjectType:
		action = rank.TagEditRank
	default:
		handler.HandleResponse(ctx, errors.BadRequest(reason.ObjectNotFound), nil)
		return
	}

	can, err := rc.rankService.CheckOperationPermission(ctx, req.UserID, action, req.ID)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	resp, err := rc.revisionListService.CheckCanUpdateRevision(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
