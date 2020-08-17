package interactions

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

var deleteResourceByIDQuery = `DELETE FROM resources WHERE id IN (SELECT r.id FROM resources r JOIN resource_relations rr ON r.id = rr.resource_id WHERE app_id = $1 AND r.id = $2)`

func (r Repository) DeleteResourceByID(resourceID, appID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return mitigate(tx, err, "Error occurred when starting the transaction", ErrCouldNotStartTx)
	}

	stmt, err := tx.Prepare(deleteResourceByIDQuery)
	if err != nil {
		return mitigate(tx, err, "Error occurred when creating the prepared statement", ErrCouldNotCreatePS)
	}

	defer stmt.Close()

	result, err := handleExec(tx, stmt, appID, resourceID)
	if err != nil {
		return err
	}

	if err = handleRowsAffected(err, result, tx); err != nil {
		return err
	}

	return handleCommit(tx)
}

func handleExec(tx *sql.Tx, stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	result, err := stmt.Exec(args...)

	if err != nil {
		return nil, mitigate(tx, err, "Error occurred when executing the prepared statement", ErrCouldNotExecStmt)
	}
	return result, nil
}

func handleRowsAffected(err error, result sql.Result, tx *sql.Tx) error {
	affected, err := result.RowsAffected()

	if err != nil {
		return mitigate(tx, err, "Error occurred while reading the rows affected", ErrCouldNotReadRowsAffected)
	}

	if affected != 1 {
		return mitigate(tx, err, "More than 1 rows were affected, rolling back", ErrMoreThanOneRowsAffected)
	}
	return nil
}

func handleCommit(tx *sql.Tx) error {
	err := tx.Commit()
	if err != nil {
		log.Errorf("Could not commit the deletion : %v", err)
		return ErrCouldNotCommit
	}
	return nil
}

func mitigate(tx *sql.Tx, err error, errMessage string, newErr error) error {
	if tx != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			log.Errorf("%s : %v", errMessage, err2)
		}
	}
	log.Errorf("Error occurred while deleting a single resource : %v", err)
	return newErr
}
