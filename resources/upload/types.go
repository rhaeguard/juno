package upload

import (
	"database/sql"
	"errors"
	"mime/multipart"
)

var (
	ErrFileCouldNotBeUploaded = errors.New("file could not be uploaded")
	EmptyID                   = ""
)

var (
	errCouldNotPersist = errors.New("could not persist the given data to database")
)

type repository interface {
	saveUploadedResourceInformation(params *SaveUploadedResourceParameters) *InsertResult
}

type FileWriter interface {
	SaveFileTo(file *multipart.FileHeader, dst string) error
}

type SaveUploadedResourceParameters struct {
	FileName          string
	FileSize          int64
	FileExtension     string
	UploadDestination string
	AppID             string
}

type InsertResult struct {
	ID  string
	Err error
}

type Repository struct {
	db *sql.DB
}

type FileUploadParameters struct {
	Hidden, Compress  bool
	Name              string
	Extension         string
	DuplicateStrategy int
	AppID             string
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
