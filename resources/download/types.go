package download

import (
	"database/sql"
	"errors"
	"time"
)

type Repository struct {
	db *sql.DB
}

type Resource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Extension string    `json:"extension"`
	Size      int64     `json:"size"`
	CreatedOn time.Time `json:"created_on"`
}

type ResourceInformation struct {
	Resources []Resource
	Err       error
}

type DownloadableResource struct {
	Resource      Resource
	SavedLocation string
}

type Service struct {
	r *Repository
}

type SingleResourceRequestParams struct {
	ResourceID, AppID, Name string
	Download                bool
}

type SingleResourceResult struct {
	File   *SingleResourceFileResult
	Data   interface{}
	Status int
}

type SingleResourceFileResult struct {
	Path, Name string
}

var (
	ErrCouldNotRetrieveResults = errors.New("could not retrieve the results")
	NoDownloadableResource     = DownloadableResource{}
)

func NewService(r *Repository) *Service {
	return &Service{r}
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}
