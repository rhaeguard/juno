package upload

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mensurowary/juno/config"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
)

func (s *Service) HandleUpload(writer FileWriter, fileHeader *multipart.FileHeader, appId string, values url.Values) (string, error) {
	parameters := makeFileUploadParams(fileHeader, values, appId)
	return s.upload(writer, fileHeader, &parameters)
}

func (s *Service) upload(writer FileWriter, file *multipart.FileHeader, params *FileUploadParameters) (string, error) {
	uploadDestination := uploadDestination(params)
	extension := fileExtension(uploadDestination)

	err := writer.SaveFileTo(file, uploadDestination)

	if err != nil {
		log.Error("Error occurred while uploading the file", params, err)
		return EmptyId, FileCouldNotBeUploaded
	}

	log.Infof("Uploaded to %s", uploadDestination)

	uploadDestination = getFilename(uploadDestination)

	res := s.r.saveUploadedResourceInformation(&SaveUploadedResourceParameters{
		FileName:          params.Name,
		FileSize:          file.Size,
		FileExtension:     extension,
		UploadDestination: uploadDestination,
		AppId:             params.AppId,
	})

	if res.Err == CouldNotPersist {
		return EmptyId, FileCouldNotBeUploaded
	}
	return res.ID, nil
}

func getFilename(uploadDestination string) string {
	start := len(config.Config.FileUploadDir) + 1
	return uploadDestination[start:]
}

func fileExtension(uploadDestination string) string {
	ext := filepath.Ext(uploadDestination)
	if ext != "" {
		return ext[1:]
	}
	return ext
}

func makeFileUploadParams(fileHeader *multipart.FileHeader, values url.Values, appId string) FileUploadParameters {
	name := strings.TrimSpace(values.Get("name"))
	if name == "" {
		name = fileHeader.Filename
	}
	return FileUploadParameters{
		Name:  name,
		AppId: appId,
	}
}

func uploadDestination(params *FileUploadParameters) string {
	ext := filepath.Ext(params.Name)
	uploadDestination := filepath.Join(config.Config.FileUploadDir, params.Name)
	uploadDestination = strings.Replace(uploadDestination, ext, "", len(uploadDestination)-len(ext))
	return fmt.Sprintf("%s-%s%s", uploadDestination, uuid.New().String(), ext)
}
