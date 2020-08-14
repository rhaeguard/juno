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
	Id        string    `json:"id"`
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

var CouldNotRetrieveResults = errors.New("could not retrieve the results")
var NoDownloadableResource = DownloadableResource{}

func NewService(db *sql.DB) Service {
	return Service{
		r: Repository{
			db: db,
		},
	}
}
