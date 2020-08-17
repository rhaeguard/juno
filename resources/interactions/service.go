package interactions

import (
	"github.com/mensurowary/juno/config"
	"github.com/mensurowary/juno/resources/download"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func (s *Service) DeleteSingleResourceByID(resourceID, appID string) error {
	resourceInfo := s.rs.GetSingleResourceInformation(download.SingleResourceRequestParams{
		ResourceID: resourceID,
		AppID:      appID,
	})

	if resourceInfo == download.NoDownloadableResource {
		log.Infof("Requested resource [%s] does not exist", resourceID)
		return ErrCouldNotFind
	}

	err := s.r.DeleteResourceByID(resourceID, appID)
	if err != nil {
		log.Errorf("Could not delete the resource [%s]", resourceID)
		return ErrCouldNotDeleteData
	}

	location := filepath.Join(config.Config.FileUploadDir, resourceInfo.SavedLocation)
	if err := os.Remove(location); err != nil {
		log.Errorf(`Error occurred while deleting the file "%s" : %v`, location, err)
		return ErrCouldNotDeleteFile
	}
	return nil
}
