package upload

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mensurowary/juno/auth"
	"github.com/mensurowary/juno/config"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strings"
)

var FileCouldNotBeUploaded = errors.New("file could not be uploaded")
var EmptyId = ""

type repository interface {
	saveUploadedResourceInformation(params *SaveUploadedResourceParameters) *InsertResult
}

type Service struct {
	r repository
}

func NewService(db *sql.DB) *Service {
	return &Service{
		r: &Repository{
			db: db,
		},
	}
}

func (s *Service) HandleUpload(c *gin.Context, fileHeader *multipart.FileHeader, values url.Values) (string, error) {
	appId := auth.GetAppId(c)
	parameters := makeFileUploadParams(fileHeader, values, appId)
	return s.upload(c, fileHeader, &parameters)
}

func (s *Service) upload(c *gin.Context, file *multipart.FileHeader, params *FileUploadParameters) (string, error) {
	uploadDestination := uploadDestination(params)
	extension := fileExtension(uploadDestination)

	err := c.SaveUploadedFile(file, uploadDestination)

	if err != nil {
		log.Error("Error occurred while uploading the file", params, err)
		return EmptyId, FileCouldNotBeUploaded
	}

	log.Infof("Uploaded to %s", uploadDestination)

	start := len(config.Config.FileUploadDir) + 1
	uploadDestination = uploadDestination[start:] // removing the directory name

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

type FileUploadParameters struct {
	Hidden, Compress  bool
	Name              string
	DuplicateStrategy int
	AppId             string
}
