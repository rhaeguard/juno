package interactions

import (
	"errors"
	"github.com/mensurowary/juno/config"
	"github.com/mensurowary/juno/resources/download"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

type Service struct {
	r               Repository
	resourceService download.Service
}

var CouldNotDelete = errors.New("could not delete the resource information from database")
var CouldNotDeleteFile = errors.New("could not delete the file")
var CouldNotFind = errors.New("could not find the resource")

func (s Service) DeleteSingleResourceById(resourceId, appId string) error {
	resourceInfo := s.resourceService.GetSingleResourceInformation(download.SingleResourceRequestParams{
		ResourceId: resourceId,
		AppId:      appId,
	})

	if resourceInfo == download.NoDownloadableResource {
		log.Infof("Requested resource [%s] does not exist", resourceId)
		return CouldNotFind
	}

	err := s.r.DeleteResourceById(resourceId, appId)
	if err != nil {
		log.Errorf("Could not delete the resource [%s]", resourceId)
		return CouldNotDelete
	}

	location := filepath.Join(config.Config.FileUploadDir, resourceInfo.SavedLocation)
	if err := os.Remove(location); err != nil {
		log.Errorf(`Error occurred while deleting the file "%s" : %v`, location, err)
		return CouldNotDeleteFile
	}
	return nil
}
