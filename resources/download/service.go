package download

import (
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/config"
	"net/http"
	"path/filepath"
	"strings"
)

func (s *Service) GetAppResourcesInformation(appID string) ResourceInformation {
	resources, err := s.r.GetResourcesByApplication(appID)
	if resources == nil {
		resources = []Resource{}
	}
	return ResourceInformation{
		Resources: resources,
		Err:       err,
	}
}

func (s *Service) GetSingleResource(params SingleResourceRequestParams) SingleResourceResult {
	downloadableResource := s.r.FindResourceLocation(params.AppID, params.ResourceID)

	if downloadableResource == NoDownloadableResource {
		return SingleResourceResult{
			File:   nil,
			Data:   commons.MakeFailureResponse("Could not find the requested resource", http.StatusNotFound),
			Status: http.StatusNotFound,
		}
	}

	if params.Download {
		path := filepath.Join(config.Config.FileUploadDir, downloadableResource.SavedLocation)
		name := getFileName(&params, &downloadableResource.Resource)
		return SingleResourceResult{
			File: &SingleResourceFileResult{
				Name: name,
				Path: path,
			},
		}
	}

	return SingleResourceResult{
		File:   nil,
		Data:   commons.MakeSuccessResponse("Successfully retrieved the resource information", downloadableResource.Resource),
		Status: http.StatusOK,
	}
}

func (s *Service) GetSingleResourceInformation(params SingleResourceRequestParams) DownloadableResource {
	return s.r.FindResourceLocation(params.AppID, params.ResourceID)
}

func getFileName(p *SingleResourceRequestParams, r *Resource) string {
	result := r.Name

	if name := strings.TrimSpace(p.Name); name != "" {
		result = name
	}

	if r.Extension != "" {
		result += "." + r.Extension
	}

	return result
}
