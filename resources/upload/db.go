package upload

import (
	"database/sql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

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
		return errCouldNotPersist
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
			return errCouldNotPersist
		}
		log.Info("Successfully inserted the data")
		return nil
	})
}

func (r *Repository) saveUploadedResourceInfo(tx *sql.Tx, ID string, params *SaveUploadedResourceParameters) error {
	return execute(tx,
		`INSERT INTO resources(id, name, extension, size, created_on) values ($1, $2, $3, $4, current_timestamp)`,
		ID, params.FileName, params.FileExtension, params.FileSize)
}

func (r *Repository) persistResourceRelations(tx *sql.Tx, resourceID string, params *SaveUploadedResourceParameters) error {
	return execute(tx,
		`INSERT INTO resource_relations(app_id, resource_id, saved_location) values ($1, $2, $3)`,
		params.AppID, resourceID, params.UploadDestination)

}

func execute(tx *sql.Tx, query string, args ...interface{}) error {
	stmt, err := prepare(tx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		log.Errorf("Could not execute the statement : %v", err)
		return errCouldNotPersist
	}

	return handleAffectedRows(tx, result)
}

func prepare(tx *sql.Tx, query string) (*sql.Stmt, error) {
	stmt, err := tx.Prepare(query)
	if err != nil {
		log.Errorf("Could not get the prepared statement : %v", err)
		return nil, errCouldNotPersist
	}
	return stmt, nil
}

func handleAffectedRows(tx *sql.Tx, result sql.Result) error {
	if affected, err := result.RowsAffected(); err != nil || affected != 1 {
		log.Errorf("Could not persist the given data to database : affected rows : %d, error : %v", affected, err)
		err := tx.Rollback()
		if err != nil {
			log.Errorf("Could not rollback! : %v", err)
		}
		return errCouldNotPersist
	}
	return nil
}
