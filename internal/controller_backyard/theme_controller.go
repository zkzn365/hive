package controller_backyard

import (
	"answer/internal/base/handler"
	"answer/internal/schema"

	"github.com/gin-gonic/gin"
)

type ThemeController struct{}

// NewThemeController new theme controller.
func NewThemeController() *ThemeController {
	return &ThemeController{}
}

// GetThemeOptions godoc
// @Summary Get theme options
// @Description Get theme options
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/theme/options [get]
func (t *ThemeController) GetThemeOptions(ctx *gin.Context) {
	handler.HandleResponse(ctx, nil, schema.GetThemeOptions)
}
