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

func (s *Service) HandleUpload(writer FileWriter, fileHeader *multipart.FileHeader, appID string, values url.Values) (string, error) {
	parameters := makeFileUploadParams(fileHeader, values, appID)
	return s.upload(writer, fileHeader, &parameters)
}

func (s *Service) upload(writer FileWriter, file *multipart.FileHeader, params *FileUploadParameters) (string, error) {
	uploadDestination := uploadDestination(params)

	err := writer.SaveFileTo(file, uploadDestination)

	if err != nil {
		log.Error("Error occurred while uploading the file", params, err)
		return EmptyID, ErrFileCouldNotBeUploaded
	}

	log.Infof("Uploaded to %s", uploadDestination)

	uploadDestination = getFilename(uploadDestination)

	res := s.r.saveUploadedResourceInformation(&SaveUploadedResourceParameters{
		FileName:          params.Name,
		FileSize:          file.Size,
		FileExtension:     params.Extension,
		UploadDestination: uploadDestination,
		AppID:             params.AppID,
	})

	if res.Err == errCouldNotPersist {
		return EmptyID, ErrFileCouldNotBeUploaded
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

func makeFileUploadParams(fileHeader *multipart.FileHeader, values url.Values, appID string) FileUploadParameters {
	name := strings.TrimSpace(values.Get("name"))

	ext := fileExtension(fileHeader.Filename)
	if ext == "" {
		ext = fileExtension(name)
	}

	if name == "" {
		filename := fileHeader.Filename
		name = filename[:len(filename)-len(ext)-1]
	}

	return FileUploadParameters{
		Name:      name,
		Extension: ext,
		AppID:     appID,
	}
}

func uploadDestination(params *FileUploadParameters) string {
	ext := params.Extension
	if ext != "" {
		ext = "." + ext
	}
	uploadDestination := filepath.Join(config.Config.FileUploadDir, params.Name)
	return fmt.Sprintf("%s-%s%s", uploadDestination, uuid.New().String(), ext)
}
