package interactions

import (
	"database/sql"
	"errors"
	"github.com/mensurowary/juno/resources/download"
)

type Repository struct {
	db *sql.DB
}

type Service struct {
	r  *Repository
	rs resourceService
}

type resourceService interface {
	GetSingleResourceInformation(params download.SingleResourceRequestParams) download.DownloadableResource
}

// Action errors
var (
	ErrCouldNotDeleteData = errors.New("could not delete the resource information from database")
	ErrCouldNotDeleteFile = errors.New("could not delete the file")
	ErrCouldNotFind       = errors.New("could not find the resource")
)

// Database action errors
var (
	ErrCouldNotStartTx          = errors.New("could not start the transaction")
	ErrCouldNotCreatePS         = errors.New("could not create the prepared statement")
	ErrCouldNotExecStmt         = errors.New("could not execute the prepared statement")
	ErrCouldNotReadRowsAffected = errors.New("could read rows affected")
	ErrMoreThanOneRowsAffected  = errors.New("could read rows affected")
	ErrCouldNotCommit           = errors.New("could not commit the deletion")
)

func NewService(r *Repository, rs resourceService) *Service {
	return &Service{
		r:  r,
		rs: rs,
	}
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}
