package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/resources/download"
	"github.com/mensurowary/juno/resources/upload"
	"github.com/mensurowary/juno/util"
	"mime/multipart"
	"net/url"
)

func Upload(handler uploadHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		UploadHandler(wc, handler)
	}
}

func GetAppResourcesInformation(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		GetAppResourcesInformationHandler(wc, handler)
	}
}

func DownloadSingleAppResource(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		DownloadSingleAppResourceHandler(wc, handler)
	}
}

func DeleteSingleAppResource(handler resourceInteractionHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		DeleteSingleAppResourceHandler(wc, handler)
	}
}

type uploadHandler interface {
	HandleUpload(writer upload.FileWriter, fileHeader *multipart.FileHeader, appId string, values url.Values) (string, error)
}

type resourcesHandler interface {
	GetAppResourcesInformation(appId string) download.ResourceInformation
	GetSingleResource(params download.SingleResourceRequestParams) download.SingleResourceResult
}

type resourceInteractionHandler interface {
	DeleteSingleResourceById(resourceId, appId string) error
}

type UploadResult struct {
	FileId string `json:"resourceId"`
}
