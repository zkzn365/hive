package uploader

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"answer/internal/base/handler"
	"answer/internal/base/reason"
	"answer/internal/service/service_config"
	"answer/internal/service/siteinfo_common"
	"answer/pkg/dir"
	"answer/pkg/uid"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/segmentfault/pacman/errors"
)

const (
	avatarSubPath      = "avatar"
	avatarThumbSubPath = "avatar_thumb"
	postSubPath        = "post"
	brandingSubPath    = "branding"
)

var (
	subPathList = []string{
		avatarSubPath,
		avatarThumbSubPath,
		postSubPath,
		brandingSubPath,
	}
	FormatExts = map[string]imaging.Format{
		".jpg":  imaging.JPEG,
		".jpeg": imaging.JPEG,
		".png":  imaging.PNG,
		".gif":  imaging.GIF,
		".tif":  imaging.TIFF,
		".tiff": imaging.TIFF,
		".bmp":  imaging.BMP,
	}
)

// UploaderService user service
type UploaderService struct {
	serviceConfig   *service_config.ServiceConfig
	siteInfoService *siteinfo_common.SiteInfoCommonService
}

// NewUploaderService new upload service
func NewUploaderService(serviceConfig *service_config.ServiceConfig,
	siteInfoService *siteinfo_common.SiteInfoCommonService) *UploaderService {
	for _, subPath := range subPathList {
		err := dir.CreateDirIfNotExist(filepath.Join(serviceConfig.UploadPath, subPath))
		if err != nil {
			panic(err)
		}
	}
	return &UploaderService{
		serviceConfig:   serviceConfig,
		siteInfoService: siteInfoService,
	}
}

// UploadAvatarFile upload avatar file
func (us *UploaderService) UploadAvatarFile(ctx *gin.Context) (url string, err error) {
	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 5*1024*1024)
	_, file, err := ctx.Request.FormFile("file")
	if err != nil {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if _, ok := FormatExts[fileExt]; !ok {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	avatarFilePath := path.Join(avatarSubPath, newFilename)
	return us.uploadFile(ctx, file, avatarFilePath)
}

func (us *UploaderService) AvatarThumbFile(ctx *gin.Context, uploadPath, fileName string, size int) (
	avatarfile []byte, err error) {
	if size > 1024 {
		size = 1024
	}
	thumbFileName := fmt.Sprintf("%d_%d@%s", size, size, fileName)
	thumbfilePath := fmt.Sprintf("%s/%s/%s", uploadPath, avatarThumbSubPath, thumbFileName)
	avatarfile, err = os.ReadFile(thumbfilePath)
	if err == nil {
		return avatarfile, nil
	}
	filePath := fmt.Sprintf("%s/avatar/%s", uploadPath, fileName)
	avatarfile, err = os.ReadFile(filePath)
	if err != nil {
		return avatarfile, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	reader := bytes.NewReader(avatarfile)
	img, err := imaging.Decode(reader)
	if err != nil {
		return avatarfile, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	new_image := imaging.Fill(img, size, size, imaging.Center, imaging.Linear)
	var buf bytes.Buffer
	fileSuffix := path.Ext(fileName)

	_, ok := FormatExts[fileSuffix]

	if !ok {
		return avatarfile, fmt.Errorf("img extension not exist")
	}
	err = imaging.Encode(&buf, new_image, FormatExts[fileSuffix])

	if err != nil {
		return avatarfile, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	thumbReader := bytes.NewReader(buf.Bytes())
	dir.CreateDirIfNotExist(path.Join(us.serviceConfig.UploadPath, avatarThumbSubPath))
	avatarFilePath := path.Join(avatarThumbSubPath, thumbFileName)
	savefilePath := path.Join(us.serviceConfig.UploadPath, avatarFilePath)
	out, err := os.Create(savefilePath)
	if err != nil {
		return avatarfile, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	defer out.Close()
	_, err = io.Copy(out, thumbReader)
	if err != nil {
		return avatarfile, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return buf.Bytes(), nil
}

func (us *UploaderService) UploadPostFile(ctx *gin.Context) (
	url string, err error) {
	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10*1024*1024)
	_, file, err := ctx.Request.FormFile("file")
	if err != nil {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if _, ok := FormatExts[fileExt]; !ok {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	avatarFilePath := path.Join(postSubPath, newFilename)
	return us.uploadFile(ctx, file, avatarFilePath)
}

func (us *UploaderService) UploadBrandingFile(ctx *gin.Context) (
	url string, err error) {
	// max size
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 10*1024*1024)
	_, file, err := ctx.Request.FormFile("file")
	if err != nil {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}
	fileExt := strings.ToLower(path.Ext(file.Filename))
	_, ok := FormatExts[fileExt]
	if !ok && fileExt != ".ico" {
		handler.HandleResponse(ctx, errors.BadRequest(reason.RequestFormatError), nil)
		return
	}

	newFilename := fmt.Sprintf("%s%s", uid.IDStr12(), fileExt)
	avatarFilePath := path.Join(brandingSubPath, newFilename)
	return us.uploadFile(ctx, file, avatarFilePath)
}

func (us *UploaderService) uploadFile(ctx *gin.Context, file *multipart.FileHeader, fileSubPath string) (
	url string, err error) {
	siteGeneral, err := us.siteInfoService.GetSiteGeneral(ctx)
	if err != nil {
		return "", err
	}
	filePath := path.Join(us.serviceConfig.UploadPath, fileSubPath)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	url = fmt.Sprintf("%s/uploads/%s", siteGeneral.SiteUrl, fileSubPath)
	return url, nil
}
