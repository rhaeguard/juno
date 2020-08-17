package resources

import (
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/resources/download"
	"github.com/mensurowary/juno/resources/interactions"
	"github.com/mensurowary/juno/resources/upload"
	"github.com/mensurowary/juno/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// UploadHandler handles the overall flow of the file uploading
func UploadHandler(wc *util.WebContext, handler uploadHandler) {
	file, err := wc.FormFile()
	if err != nil {
		log.Errorf("Error occurred while retrieving the file from request : %s", err)
		wc.UnprocessableEntity(commons.MakeFailureResponse(
			"Could not retrieve the uploaded file from the request", http.StatusUnprocessableEntity,
		))
		return
	}

	appID := wc.GetAppID()

	ID, err := handler.HandleUpload(wc, file, appID, wc.Form())
	if err == upload.ErrFileCouldNotBeUploaded || ID == upload.EmptyID {
		wc.UnprocessableEntity(commons.MakeFailureResponse(
			"File could not be uploaded", http.StatusUnprocessableEntity,
		))
	} else {
		wc.Ok(commons.MakeSuccessResponse(
			"Successfully uploaded the file",
			UploadResult{
				FileID: ID,
			}),
		)
	}
}

func DeleteSingleAppResourceHandler(wc *util.WebContext, handler resourceInteractionHandler) {
	resourceID := wc.GetResourceID()
	appID := wc.GetAppID()
	if err := handler.DeleteSingleResourceByID(resourceID, appID); err != nil {
		switch err {
		case interactions.ErrCouldNotDeleteData:
			wc.UnprocessableEntity(commons.MakeFailureResponse("Could not delete the resource information", http.StatusUnprocessableEntity))
		case interactions.ErrCouldNotDeleteFile:
			wc.UnprocessableEntity(commons.MakeFailureResponse("Could not delete the resource file", http.StatusUnprocessableEntity))
		case interactions.ErrCouldNotFind:
			wc.NotFound(commons.MakeFailureResponse("Could not find the requested resource", http.StatusNotFound))
		default:
			wc.InternalServerError(commons.MakeFailureResponse("Unknown error occurred", http.StatusInternalServerError))
		}
	} else {
		wc.Ok(commons.MakeSuccessResponse("Successfully deleted the resource", nil))
	}
}

func GetAppResourcesInformationHandler(wc *util.WebContext, handler resourcesHandler) {
	appID := wc.GetAppID()
	if info := handler.GetAppResourcesInformation(appID); info.Err != nil {
		wc.NotFound(commons.MakeFailureResponse("Could not retrieve the data", http.StatusNotFound))
	} else {
		wc.Ok(commons.MakeSuccessResponse("Successfully retrieved all the available resources", info.Resources))
	}
}

func DownloadSingleAppResourceHandler(wc *util.WebContext, handler resourcesHandler) {
	params := getSingleResourceParams(wc)
	result := handler.GetSingleResource(params)
	if result.File != nil {
		wc.RespondWithFile(result.File.Path, result.File.Name)
	} else {
		wc.Respond(result.Status, result.Data)
	}
}

func getSingleResourceParams(wc *util.WebContext) download.SingleResourceRequestParams {
	name := wc.QueryParam("name")
	downloadParam := wc.QueryParam("download")

	shouldDownload := strings.ToLower(downloadParam) == "true"

	params := download.SingleResourceRequestParams{
		ResourceID: wc.GetResourceID(),
		AppID:      wc.GetAppID(),
		Name:       name,
		Download:   shouldDownload,
	}
	return params
}
