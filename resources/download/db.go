package download

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

func (r *Repository) GetResourcesByApplication(appID string) ([]Resource, error) {
	rows, err := r.queryAllForAppID(appID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		id, name, extension string
		size                int64
		createdOn           time.Time
	)
	var resources []Resource
	for rows.Next() {
		err := rows.Scan(&id, &name, &extension, &size, &createdOn)
		if err != nil {
			log.Errorf("Error occurred while mapping the results to objects : %v", err)
			return nil, ErrCouldNotRetrieveResults
		}
		resources = append(resources, Resource{
			ID:        id,
			Name:      name,
			Extension: extension,
			Size:      size,
			CreatedOn: createdOn,
		})
	}
	return resources, nil
}

func (r *Repository) queryAllForAppID(appID string) (*sql.Rows, error) {
	rows, err := r.db.Query(`SELECT r.id, r.name, r.extension, r.size, r.created_on FROM resources r JOIN resource_relations rr ON r.id = rr.resource_id WHERE rr.app_id = $1`, appID)

	if err != nil {
		log.Errorf("Error occurred while trying to retrieve resources for the app: %s : %v", appID, err)
		return nil, ErrCouldNotRetrieveResults
	}
	return rows, nil
}

func (r *Repository) FindResourceLocation(appID, resourceID string) DownloadableResource {
	var (
		name, extension, savedLocation, id string
		size                               int64
		createdOn                          time.Time
	)

	rows := r.queryForResourceInformation(appID, resourceID)
	if err := rows.Scan(&id, &name, &extension, &size, &createdOn, &savedLocation); err != nil {
		log.Errorf("Error occurred while mapping the results to objects : %v", err)
		return NoDownloadableResource
	}
	return DownloadableResource{
		Resource: Resource{
			ID:        id,
			Name:      name,
			Extension: extension,
			CreatedOn: createdOn,
			Size:      size,
		},
		SavedLocation: savedLocation,
	}
}

func (r *Repository) queryForResourceInformation(appID string, resourceID string) *sql.Row {
	return r.db.QueryRow(`
		SELECT r.id, r.name, r.extension, r.size, r.created_on, rr.saved_location
		FROM resources r 
		JOIN resource_relations rr ON r.id = rr.resource_id
		WHERE rr.app_id = $1 AND r.id = $2
	`, appID, resourceID)
}
