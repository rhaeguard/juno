package upload

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var (
	CouldNotPersist = errors.New("could not persist the given data to database")
)

type SaveUploadedResourceParameters struct {
	FileName          string
	FileSize          int64
	FileExtension     string
	UploadDestination string
	AppId             string
}

type InsertResult struct {
	ID  string
	Err error
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) saveUploadedResourceInformation(params *SaveUploadedResourceParameters) *InsertResult {
	ID := uuid.New().String()
	err := r.persistResourceInformation(ID, params)
	if err != nil {
		ID = ""
	}
	return &InsertResult{
		ID:  ID,
		Err: err,
	}
}

func (r *Repository) withTx(action func(tx *sql.Tx) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Errorf("Could not start the transaction : %v", err)
		return CouldNotPersist
	}
	return action(tx)
}

func (r *Repository) persistResourceInformation(ID string, params *SaveUploadedResourceParameters) error {
	return r.withTx(func(tx *sql.Tx) error {
		if err := r.saveUploadedResourceInfo(tx, ID, params); err != nil {
			return err
		}

		if err := r.persistResourceRelations(tx, ID, params); err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			log.Infof("Could not commit! : %v", err)
			return CouldNotPersist
		}
		log.Info("Successfully inserted the data")
		return nil
	})
}

func (r *Repository) saveUploadedResourceInfo(tx *sql.Tx, ID string, params *SaveUploadedResourceParameters) error {
	stmt, err := tx.Prepare(`
			INSERT INTO resources(id, name, extension, size, created_on) 
			values ($1, $2, $3, $4, current_timestamp)
	`)
	if err != nil {
		log.Errorf("Could not get the prepared statement : %v", err)
		return CouldNotPersist
	}

	defer stmt.Close()

	result, err := stmt.Exec(ID, params.FileName, params.FileExtension, params.FileSize)
	if err != nil {
		log.Errorf("Could not execute the statement : %v", err)
		return CouldNotPersist
	}

	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		log.Errorf("Could not persist the given data to database : affected rows : %d, error : %v", affected, err)
		err := tx.Rollback()
		if err != nil {
			log.Errorf("Could not rollback : %v", err)
		}
		return CouldNotPersist
	}
	return nil
}

func (r *Repository) persistResourceRelations(tx *sql.Tx, resourceId string, params *SaveUploadedResourceParameters) error {
	stmt, err := tx.Prepare(`
			INSERT INTO resource_relations(app_id, resource_id, saved_location) 
			values ($1, $2, $3)
	`)
	if err != nil {
		log.Errorf("Could not get the prepared statement : %v", err)
		return CouldNotPersist
	}

	result, err := stmt.Exec(params.AppId, resourceId, params.UploadDestination)

	if err != nil {
		log.Errorf("Could not execute the statement : %v", err)
		return CouldNotPersist
	}

	if affected, err := result.RowsAffected(); err != nil || affected != 1 {
		log.Info("Could not insert")
		err := tx.Rollback()
		if err != nil {
			log.Errorf("Could not rollback! : %v", err)
		}
		return CouldNotPersist
	}
	return nil
}
