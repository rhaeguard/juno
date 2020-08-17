package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/resources/download"
	"github.com/mensurowary/juno/resources/upload"
	"github.com/mensurowary/juno/util"
	"mime/multipart"
	"net/url"
)

// Upload is the upload handler
// Deals with uploading a file and its different modes
func Upload(handler uploadHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		UploadHandler(wc, handler)
	}
}

// GetAppResourcesInformation retrieves the resources information
func GetAppResourcesInformation(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		GetAppResourcesInformationHandler(wc, handler)
	}
}

// DownloadSingleAppResource handles single resource information retrieval/file download
func DownloadSingleAppResource(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		DownloadSingleAppResourceHandler(wc, handler)
	}
}

// DeleteSingleAppResource handles the deletion of a resource
func DeleteSingleAppResource(handler resourceInteractionHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		wc := util.NewWebContext(c)
		DeleteSingleAppResourceHandler(wc, handler)
	}
}

type uploadHandler interface {
	HandleUpload(writer upload.FileWriter, fileHeader *multipart.FileHeader, appID string, values url.Values) (string, error)
}

type resourcesHandler interface {
	GetAppResourcesInformation(appID string) download.ResourceInformation
	GetSingleResource(params download.SingleResourceRequestParams) download.SingleResourceResult
}

type resourceInteractionHandler interface {
	DeleteSingleResourceByID(resourceID, appID string) error
}

// UploadResult represents the result of the file upload
type UploadResult struct {
	FileID string `json:"resourceId"`
}
