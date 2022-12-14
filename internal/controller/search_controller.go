package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/schema"
	"answer/internal/service"

	"github.com/gin-gonic/gin"
)

// SearchController tag controller
type SearchController struct {
	searchService *service.SearchService
}

// NewSearchController new controller
func NewSearchController(searchService *service.SearchService) *SearchController {
	return &SearchController{searchService: searchService}
}

// Search godoc
// @Summary search object
// @Description search object
// @Tags Search
// @Produce json
// @Security ApiKeyAuth
// @Param q query string true "query string"
// @Param order query string true "order" Enums(newest,active,score,relevance)
// @Success 200 {object} handler.RespBody{data=schema.SearchListResp}
// @Router /answer/api/v1/search [get]
func (sc *SearchController) Search(ctx *gin.Context) {
	dto := schema.SearchDTO{}

	if handler.BindAndCheck(ctx, &dto) {
		return
	}
	dto.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, total, extra, err := sc.searchService.Search(ctx, &dto)

	handler.HandleResponse(ctx, err, schema.SearchListResp{
		Total:      total,
		SearchResp: resp,
		Extra:      extra,
	})
}
