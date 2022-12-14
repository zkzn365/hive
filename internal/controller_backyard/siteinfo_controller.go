package controller_backyard

import (
	"answer/internal/base/handler"
	"answer/internal/base/middleware"
	"answer/internal/schema"
	"answer/internal/service/siteinfo"

	"github.com/gin-gonic/gin"
)

// SiteInfoController site info controller
type SiteInfoController struct {
	siteInfoService *siteinfo.SiteInfoService
}

// NewSiteInfoController new site info controller
func NewSiteInfoController(siteInfoService *siteinfo.SiteInfoService) *SiteInfoController {
	return &SiteInfoController{
		siteInfoService: siteInfoService,
	}
}

// GetGeneral get site general information
// @Summary get site general information
// @Description get site general information
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.SiteGeneralResp}
// @Router /answer/admin/api/siteinfo/general [get]
func (sc *SiteInfoController) GetGeneral(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSiteGeneral(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// GetInterface get site interface
// @Summary get site interface
// @Description get site interface
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.SiteInterfaceResp}
// @Router /answer/admin/api/siteinfo/interface [get]
func (sc *SiteInfoController) GetInterface(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSiteInterface(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// GetSiteBranding get site interface
// @Summary get site interface
// @Description get site interface
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.SiteBrandingResp}
// @Router /answer/admin/api/siteinfo/branding [get]
func (sc *SiteInfoController) GetSiteBranding(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSiteBranding(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// GetSiteWrite get site interface
// @Summary get site interface
// @Description get site interface
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.SiteWriteResp}
// @Router /answer/admin/api/siteinfo/write [get]
func (sc *SiteInfoController) GetSiteWrite(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSiteWrite(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// GetSiteLegal Set the legal information for the site
// @Summary Set the legal information for the site
// @Description Set the legal information for the site
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.SiteLegalResp}
// @Router /answer/admin/api/siteinfo/legal [get]
func (sc *SiteInfoController) GetSiteLegal(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSiteLegal(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// UpdateGeneral update site general information
// @Summary update site general information
// @Description update site general information
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.SiteGeneralReq true "general"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/siteinfo/general [put]
func (sc *SiteInfoController) UpdateGeneral(ctx *gin.Context) {
	req := schema.SiteGeneralReq{}
	if handler.BindAndCheck(ctx, &req) {
		return
	}
	err := sc.siteInfoService.SaveSiteGeneral(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateInterface update site interface
// @Summary update site info interface
// @Description update site info interface
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.SiteInterfaceReq true "general"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/siteinfo/interface [put]
func (sc *SiteInfoController) UpdateInterface(ctx *gin.Context) {
	req := schema.SiteInterfaceReq{}
	if handler.BindAndCheck(ctx, &req) {
		return
	}
	err := sc.siteInfoService.SaveSiteInterface(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateBranding update site branding
// @Summary update site info branding
// @Description update site info branding
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.SiteBrandingReq true "branding info"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/siteinfo/branding [put]
func (sc *SiteInfoController) UpdateBranding(ctx *gin.Context) {
	req := &schema.SiteBrandingReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := sc.siteInfoService.SaveSiteBranding(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// UpdateSiteWrite update site write info
// @Summary update site write info
// @Description update site write info
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.SiteWriteReq true "write info"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/siteinfo/write [put]
func (sc *SiteInfoController) UpdateSiteWrite(ctx *gin.Context) {
	req := &schema.SiteWriteReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	req.UserID = middleware.GetLoginUserIDFromContext(ctx)

	resp, err := sc.siteInfoService.SaveSiteWrite(ctx, req)
	handler.HandleResponse(ctx, err, resp)
}

// UpdateSiteLegal update site legal info
// @Summary update site legal info
// @Description update site legal info
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.SiteLegalReq true "write info"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/siteinfo/legal [put]
func (sc *SiteInfoController) UpdateSiteLegal(ctx *gin.Context) {
	req := &schema.SiteLegalReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := sc.siteInfoService.SaveSiteLegal(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}

// GetSMTPConfig get smtp config
// @Summary GetSMTPConfig get smtp config
// @Description GetSMTPConfig get smtp config
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Success 200 {object} handler.RespBody{data=schema.GetSMTPConfigResp}
// @Router /answer/admin/api/setting/smtp [get]
func (sc *SiteInfoController) GetSMTPConfig(ctx *gin.Context) {
	resp, err := sc.siteInfoService.GetSMTPConfig(ctx)
	handler.HandleResponse(ctx, err, resp)
}

// UpdateSMTPConfig update smtp config
// @Summary update smtp config
// @Description update smtp config
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param data body schema.UpdateSMTPConfigReq true "smtp config"
// @Success 200 {object} handler.RespBody{}
// @Router /answer/admin/api/setting/smtp [put]
func (sc *SiteInfoController) UpdateSMTPConfig(ctx *gin.Context) {
	req := &schema.UpdateSMTPConfigReq{}
	if handler.BindAndCheck(ctx, req) {
		return
	}
	err := sc.siteInfoService.UpdateSMTPConfig(ctx, req)
	handler.HandleResponse(ctx, err, nil)
}
