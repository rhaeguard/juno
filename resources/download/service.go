package download

import (
	"github.com/gin-gonic/gin"
	"github.com/mensurowary/juno/auth"
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/config"
	"net/http"
	"path/filepath"
	"strings"
)

type Service struct {
	r Repository
}

type SingleResourceRequestParams struct {
	ResourceId, AppId, Name string
	Download                bool
}

func (s Service) GetAppResourcesInformation(c *gin.Context) ResourceInformation {
	appId := auth.GetAppId(c)
	resources, err := s.r.GetResourcesByApplication(appId)
	if resources == nil {
		resources = []Resource{}
	}
	return ResourceInformation{
		Resources: resources,
		Err:       err,
	}
}

func (s Service) GetSingleResource(c *gin.Context, params SingleResourceRequestParams) {
	downloadableResource := s.r.FindResourceLocation(params.AppId, params.ResourceId)

	if downloadableResource == NoDownloadableResource {
		c.JSON(http.StatusNotFound, commons.MakeFailureResponse("Could not find the requested resource", http.StatusNotFound))
		return
	}

	if params.Download {
		filePath := filepath.Join(config.Config.FileUploadDir, downloadableResource.SavedLocation)
		c.FileAttachment(filePath, fullPath(&params, &downloadableResource.Resource))
	} else {
		c.JSON(http.StatusOK, commons.MakeSuccessResponse("Successfully retrieved the resource information", downloadableResource.Resource))
	}
}

func (s Service) GetSingleResourceInformation(params SingleResourceRequestParams) DownloadableResource {
	return s.r.FindResourceLocation(params.AppId, params.ResourceId)
}

func fullPath(p *SingleResourceRequestParams, r *Resource) string {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		if r.Extension == "" {
			return r.Name
		}
		return r.Name + "." + r.Extension
	}
	return name + "." + r.Extension
}
