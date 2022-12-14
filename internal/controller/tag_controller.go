package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/base/reason"
	"answer/internal/schema"
	"answer/internal/service/rank"
	"answer/internal/service/tag"
	"answer/internal/service/tag_common"

	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

// TagController tag controller
type TagController struct {
	tagService       *tag.TagService
	tagCommonService *tag_common.TagCommonService
	rankService      *rank.RankService
}

// NewTagController new controller
func NewTagController(
	tagService *tag.TagService,
	tagCommonService *tag_common.TagCommonService,
	rankService *rank.RankService,
) *TagController {
	return &TagController{tagService: tagService, tagCommonService: tagCommonService, rankService: rankService}
}

// SearchTagLike get tag list
// @Summary get tag list
// @Description get tag list
// @Tags Tag
// @Produce json
// @Security ApiKeyAuth
// @Param tag query string false "tag"
// @Success 200 {object} handler.RespBody{data=[]schema.GetTagResp}
// @Router /answer/api/v1/question/tags [get]
func (tc *TagController) SearchTagLike(ctx *gin.Context) {
	req := &schema.SearchTagLikeReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.IsAdmin = middleware.GetIsAdminFromContext(ctx)
	resp, err := tc.tagCommonService.SearchTagLike(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// RemoveTag delete tag
// @Summary delete tag
// @Description delete tag
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.RemoveTagReq true "tag"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag [delete]
func (tc *TagController) RemoveTag(ctx *gin.Context) {
	req := &schema.RemoveTagReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, rank.TagDeleteRank, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagService.RemoveTag(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateTag update tag
// @Summary update tag
// @Description update tag
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.UpdateTagReq true "tag"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag [put]
func (tc *TagController) UpdateTag(ctx *gin.Context) {
	req := &schema.UpdateTagReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := tc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		rank.TagEditRank,
		rank.TagEditWithoutReviewRank,
	}, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !canList[0] {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}
	req.NoNeedReview = canList[1]

	err = tc.tagService.UpdateTag(ctx, req)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
	} else {
		handler.HandleResponse(ctx, err, &schema.UpdateTagResp{WaitForReview: !req.NoNeedReview})
	}
}

// GetTagInfo get tag one
// @Summary get tag one
// @Description get tag one
// @Tags Tag
// @Accept json
// @Produce json
// @Param tag_id query string true "tag id"
// @Param tag_name query string true "tag name"
// @Success 200 {object} handler.RespBody{data=schema.GetTagResp}
// @Router /answer/api/v1/tag [get]
func (tc *TagController) GetTagInfo(ctx *gin.Context) {
	req := &schema.GetTagInfoReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := tc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		rank.TagEditRank,
		rank.TagDeleteRank,
	}, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanEdit = canList[0]
	req.CanDelete = canList[1]

	resp, err := tc.tagService.GetTagInfo(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetTagWithPage get tag page
// @Summary get tag page
// @Description get tag page
// @Tags Tag
// @Produce json
// @Param page query int false "page size"
// @Param page_size query int false "page size"
// @Param slug_name query string false "slug_name"
// @Param query_cond query string false "query condition" Enums(popular, name, newest)
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetTagPageResp}}
// @Router /answer/api/v1/tags/page [get]
func (tc *TagController) GetTagWithPage(ctx *gin.Context) {
	req := &schema.GetTagWithPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := tc.tagService.GetTagWithPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// GetFollowingTags get following tag list
// @Summary get following tag list
// @Description get following tag list
// @Security ApiKeyAuth
// @Tags Tag
// @Produce json
// @Success 200 {object} handler.RespBody{data=[]schema.GetFollowingTagsResp}
// @Router /answer/api/v1/tags/following [get]
func (tc *TagController) GetFollowingTags(ctx *gin.Context) {
	userID := middleware.GetLoginUserIDFromContext(ctx)
	resp, err := tc.tagService.GetFollowingTags(ctx, userID)
	handler.HandleResponse(ctx, err, resp)
}

// GetTagSynonyms get tag synonyms
// @Summary get tag synonyms
// @Description get tag synonyms
// @Tags Tag
// @Produce json
// @Param tag_id query int true "tag id"
// @Success 200 {object} handler.RespBody{data=schema.GetTagSynonymsResp}
// @Router /answer/api/v1/tag/synonyms [get]
func (tc *TagController) GetTagSynonyms(ctx *gin.Context) {
	req := &schema.GetTagSynonymsReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	canList, err := tc.rankService.CheckOperationPermissions(ctx, req.UserID, []string{
		rank.TagSynonymRank,
	}, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	req.CanEdit = canList[0]

	resp, err := tc.tagService.GetTagSynonyms(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// UpdateTagSynonym update tag
// @Summary update tag
// @Description update tag
// @Tags Tag
// @Accept json
// @Produce json
// @Param data body schema.UpdateTagSynonymReq true "tag"
// @Success 200 {object} handler.RespBody
// @Router /answer/api/v1/tag/synonym [put]
func (tc *TagController) UpdateTagSynonym(ctx *gin.Context) {
	req := &schema.UpdateTagSynonymReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)
	can, err := tc.rankService.CheckOperationPermission(ctx, req.UserID, rank.TagSynonymRank, "")
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	if !can {
		handler.HandleResponse(ctx, errors.Forbidden(reason.RankFailToMeetTheCondition), nil)
		return
	}

	err = tc.tagService.UpdateTagSynonym(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
