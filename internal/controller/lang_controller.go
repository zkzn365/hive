package controller

import (
	"encoding/json"

	"answer/internal/base/handler"
	"answer/internal/base/translator"
	"answer/internal/service/siteinfo_common"

	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/i18n"
)

type LangController struct {
	translator      i18n.Translator
	siteInfoService *siteinfo_common.SiteInfoCommonService
}

// NewLangController new language controller.
func NewLangController(tr i18n.Translator, siteInfoService *siteinfo_common.SiteInfoCommonService) *LangController {
	return &LangController{translator: tr, siteInfoService: siteInfoService}
}

// GetLangMapping get language config mapping
// @Summary get language config mapping
// @Description get language config mapping
// @Tags Lang
// @Param Accept-Language header string true "Accept-Language"
// @Produce json
// @Success 200 {object} handler.RespBody{}
// @Router /answer/api/v1/language/config [get]
func (u *LangController) GetLangMapping(ctx *gin.Context) {
	data, _ := u.translator.Dump(handler.GetLang(ctx))
	var resp map[string]any
	_ = json.Unmarshal(data, &resp)
	handler.HandleResponse(ctx, nil, resp)
}

// GetAdminLangOptions Get language options
// @Summary Get language options
// @Description Get language options
// @Tags Lang
// @Produce json
// @Success 200 {object} handler.RespBody{}
// @Router /answer/api/v1/language/options [get]
// @Router /answer/admin/api/language/options [get]
func (u *LangController) GetAdminLangOptions(ctx *gin.Context) {
	handler.HandleResponse(ctx, nil, translator.LanguageOptions)
}

// GetUserLangOptions Get language options
// @Summary Get language options
// @Description Get language options
// @Tags Lang
// @Produce json
// @Success 200 {object} handler.RespBody{}
// @Router /answer/api/v1/language/options [get]
func (u *LangController) GetUserLangOptions(ctx *gin.Context) {
	siteInterfaceResp, err := u.siteInfoService.GetSiteInterface(ctx)
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}

	options := translator.LanguageOptions
	if len(siteInterfaceResp.Language) > 0 {
		defaultOption := []*translator.LangOption{
			{Label: translator.DefaultLangOption, Value: translator.DefaultLangOption},
		}
		options = append(defaultOption, options...)
	}
	handler.HandleResponse(ctx, nil, options)
}
