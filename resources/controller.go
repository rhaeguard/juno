package resources

import (
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/auth"
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/resources/download"
	"github.com/mensurowary/juno/resources/interactions"
	"github.com/mensurowary/juno/resources/upload"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type uploadHandler interface {
	HandleUpload(c *gin.Context, fileHeader *multipart.FileHeader, values url.Values) (string, error)
}

type UploadResult struct {
	FileId string `json:"resourceId"`
}

func Upload(handler uploadHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Errorf("Error occurred while retrieving the file from request : %s", err)
			c.JSON(http.StatusUnprocessableEntity, commons.MakeFailureResponse(
				"Could not retrieve the uploaded file from the request", http.StatusUnprocessableEntity,
			))
			return
		}

		ID, err := handler.HandleUpload(c, file, c.Request.Form)
		if err == upload.FileCouldNotBeUploaded || ID == upload.EmptyId {
			c.JSON(http.StatusUnprocessableEntity, commons.MakeFailureResponse(
				"File could not be uploaded", http.StatusUnprocessableEntity,
			))
		} else {
			c.JSON(http.StatusOK, commons.MakeSuccessResponse(
				"Successfully uploaded the file",
				UploadResult{
					FileId: ID,
				}),
			)
		}
	}
}

type resourcesHandler interface {
	GetAppResourcesInformation(c *gin.Context) download.ResourceInformation
	GetSingleResource(c *gin.Context, params download.SingleResourceRequestParams)
}

type resourceInteractionHandler interface {
	DeleteSingleResourceById(resourceId, appId string) error
}

func GetAppResourcesInformation(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		if info := handler.GetAppResourcesInformation(c); info.Err != nil {
			c.JSON(http.StatusNotFound, commons.MakeFailureResponse("Could not retrieve the data", http.StatusNotFound))
		} else {
			c.JSON(http.StatusOK, commons.MakeSuccessResponse("Successfully retrieved all the available resources", info.Resources))
		}
	}
}

func DownloadSingleAppResource(handler resourcesHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		params := getSingleResourceParams(c)
		handler.GetSingleResource(c, params)
	}
}

func getSingleResourceParams(c *gin.Context) download.SingleResourceRequestParams {
	query := c.Request.URL.Query()

	var name, downloadParam string

	if len(query["name"]) == 0 {
		name = ""
	} else {
		name = query["name"][0]
	}

	if len(query["download"]) == 0 {
		downloadParam = ""
	} else {
		downloadParam = query["download"][0]
	}

	shouldDownload := strings.ToLower(downloadParam) == "true"

	params := download.SingleResourceRequestParams{
		ResourceId: c.Param("id"),
		AppId:      auth.GetAppId(c),
		Name:       name,
		Download:   shouldDownload,
	}
	return params
}

func DeleteSingleAppResource(handler resourceInteractionHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		resourceId := c.Param("id")
		appId := auth.GetAppId(c)
		if err := handler.DeleteSingleResourceById(resourceId, appId); err != nil {
			switch err {
			case interactions.CouldNotDelete:
				c.JSON(http.StatusUnprocessableEntity, commons.MakeFailureResponse("Could not delete the resource information", http.StatusUnprocessableEntity))
			case interactions.CouldNotDeleteFile:
				c.JSON(http.StatusUnprocessableEntity, commons.MakeFailureResponse("Could not delete the resource file", http.StatusUnprocessableEntity))
			case interactions.CouldNotFind:
				c.JSON(http.StatusNotFound, commons.MakeFailureResponse("Could not find the requested resource", http.StatusNotFound))
			default:
				c.JSON(http.StatusInternalServerError, commons.MakeFailureResponse("Unknown error occurred", http.StatusInternalServerError))
			}
		} else {
			c.JSON(http.StatusOK, commons.MakeSuccessResponse("Successfully deleted the resource", nil))
		}

	}
}
