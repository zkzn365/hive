package controller

import (
	"answer/internal/base/handler"
	"answer/internal/base/reason"
	"answer/internal/service/uploader"

	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

const (
	// file is uploaded by markdown(or something else) editor
	fileFromPost = "post"
	// file is used to change the user's avatar
	fileFromAvatar = "avatar"
	// file is logo/icon images
	fileFromBranding = "branding"
)

// UploadController upload controller
type UploadController struct {
	uploaderService *uploader.UploaderService
}

// NewUploadController new controller
func NewUploadController(uploaderService *uploader.UploaderService) *UploadController {
	return &UploadController{
		uploaderService: uploaderService,
	}
}

// UploadFile upload file
// @Summary upload file
// @Description upload file
// @Tags Upload
// @Accept multipart/form-data
// @Security ApiKeyAuth
// @Param source formData string true "identify the source of the file upload" Enums(post, avatar, branding)
// @Param file formData file true "file"
// @Success 200 {object} handler.RespBody{data=string}
// @Router /answer/api/v1/file [post]
func (uc *UploadController) UploadFile(ctx *gin.Context) {
	var (
		url string
		err error
	)

	source := ctx.PostForm("source")
	switch source {
	case fileFromAvatar:
		url, err = uc.uploaderService.UploadAvatarFile(ctx)
	case fileFromPost:
		url, err = uc.uploaderService.UploadPostFile(ctx)
	case fileFromBranding:
		url, err = uc.uploaderService.UploadBrandingFile(ctx)
	default:
		handler.HandleResponse(ctx, errors.BadRequest(reason.UploadFileSourceUnsupported), nil)
		return
	}
	if err != nil {
		handler.HandleResponse(ctx, err, nil)
		return
	}
	handler.HandleResponse(ctx, err, url)
}
