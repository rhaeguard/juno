package interactions

import (
	"database/sql"
	"errors"
	log "github.com/sirupsen/logrus"
)

type Repository struct {
	db *sql.DB
}

var (
	ErrCouldNotStartTx          = errors.New("could not start the transaction")
	ErrCouldNotCreatePS         = errors.New("could not create the prepared statement")
	ErrCouldNotExecStmt         = errors.New("could not execute the prepared statement")
	ErrCouldNotReadRowsAffected = errors.New("could read rows affected")
	ErrMoreThanOneRowsAffected  = errors.New("could read rows affected")
	ErrCouldNotCommit           = errors.New("could not commit the deletion")
)

func (r Repository) DeleteResourceById(resourceId, appId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		mitigate(tx, err, "Error occurred when starting the transaction")
		return ErrCouldNotStartTx
	}
	stmt, err := tx.Prepare(`
		DELETE FROM resources
		WHERE id IN (
			SELECT r.id FROM resources r
			JOIN resource_relations rr ON r.id = rr.resource_id
			WHERE app_id = $1 AND r.id = $2
		)
	`)

	defer stmt.Close()

	if err != nil {
		mitigate(tx, err, "Error occurred when creating the prepared statement")
		return ErrCouldNotCreatePS
	}

	result, err := stmt.Exec(appId, resourceId)

	if err != nil {
		mitigate(tx, err, "Error occurred when executing the prepared statement")
		return ErrCouldNotExecStmt
	}

	affected, err := result.RowsAffected()

	if err != nil {
		mitigate(tx, err, "Error occurred while reading the rows affected")
		return ErrCouldNotReadRowsAffected
	}

	if affected == 1 {
		err := tx.Commit()
		if err != nil {
			log.Errorf("Could not commit the deletion : %v", err)
			return ErrCouldNotCommit
		}
		return nil
	} else {
		mitigate(tx, err, "More than 1 rows were affected, rolling back")
		return ErrMoreThanOneRowsAffected
	}
}

func mitigate(tx *sql.Tx, err error, errMessage string) {
	if tx != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Errorf("%s : %v", errMessage, err2)
		}
	}
	log.Errorf("Error occurred while deleting a single resource : %v", err)
}
