package interactions

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

func (r Repository) DeleteResourceById(resourceId, appId string) error {
	tx, err := r.db.Begin()
	if err != nil {
		mitigate(tx, err, "Error occurred when starting the transaction")
		return ErrCouldNotStartTx
	}
	stmt, err := tx.Prepare(`DELETE FROM resources WHERE id IN (SELECT r.id FROM resources r JOIN resource_relations rr ON r.id = rr.resource_id WHERE app_id = $1 AND r.id = $2)`)

	if err != nil {
		mitigate(tx, err, "Error occurred when creating the prepared statement")
		return ErrCouldNotCreatePS
	}

	defer stmt.Close()

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

	if affected != 1 {
		mitigate(tx, err, "More than 1 rows were affected, rolling back")
		return ErrMoreThanOneRowsAffected
	}

	err = tx.Commit()
	if err != nil {
		log.Errorf("Could not commit the deletion : %v", err)
		return ErrCouldNotCommit
	}
	return nil
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
