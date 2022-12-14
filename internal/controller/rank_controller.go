package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/schema"
	"answer/internal/service/rank"

	"github.com/gin-gonic/gin"
)

// RankController rank controller
type RankController struct {
	rankService *rank.RankService
}

// NewRankController new controller
func NewRankController(
	rankService *rank.RankService) *RankController {
	return &RankController{rankService: rankService}
}

// GetRankPersonalWithPage user personal rank list
// @Summary user personal rank list
// @Description user personal rank list
// @Tags Rank
// @Produce json
// @Param page query int false "page"
// @Param page_size query int false "page size"
// @Param username query string false "username"
// @Success 200 {object} handler.RespBody{data=pager.PageModel{list=[]schema.GetRankPersonalWithPageResp}}
// @Router /answer/api/v1/personal/rank/page [get]
func (cc *RankController) GetRankPersonalWithPage(ctx *gin.Context) {
	req := &schema.GetRankPersonalWithPageReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}

	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := cc.rankService.GetRankPersonalWithPage(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}
