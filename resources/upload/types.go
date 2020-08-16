package upload

import (
	"database/sql"
	"errors"
	"mime/multipart"
)

var FileCouldNotBeUploaded = errors.New("file could not be uploaded")
var EmptyId = ""

type repository interface {
	saveUploadedResourceInformation(params *SaveUploadedResourceParameters) *InsertResult
}

type FileWriter interface {
	SaveFileTo(file *multipart.FileHeader, dst string) error
}

type FileUploadParameters struct {
	Hidden, Compress  bool
	Name              string
	DuplicateStrategy int
	AppId             string
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
