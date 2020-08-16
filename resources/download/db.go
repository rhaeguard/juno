package download

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
	"time"
)

func (r *Repository) GetResourcesByApplication(appId string) ([]Resource, error) {
	rows, err := r.queryAllForAppId(appId)
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
			return nil, CouldNotRetrieveResults
		}
		resources = append(resources, Resource{
			Id:        id,
			Name:      name,
			Extension: extension,
			Size:      size,
			CreatedOn: createdOn,
		})
	}
	return resources, nil
}

func (r *Repository) queryAllForAppId(appId string) (*sql.Rows, error) {
	rows, err := r.db.Query(`SELECT r.id, r.name, r.extension, r.size, r.created_on FROM resources r JOIN resource_relations rr ON r.id = rr.resource_id WHERE rr.app_id = $1`, appId)

	if err != nil {
		log.Errorf("Error occurred while trying to retrieve resources for the app: %s : %v", appId, err)
		return nil, CouldNotRetrieveResults
	}
	return rows, nil
}

func (r *Repository) FindResourceLocation(appId, resourceId string) DownloadableResource {
	var (
		name, extension, savedLocation, id string
		size                               int64
		createdOn                          time.Time
	)

	rows := r.queryForResourceInformation(appId, resourceId)
	if err := rows.Scan(&id, &name, &extension, &size, &createdOn, &savedLocation); err != nil {
		log.Errorf("Error occurred while mapping the results to objects : %v", err)
		return NoDownloadableResource
	}
	return DownloadableResource{
		Resource: Resource{
			Id:        id,
			Name:      name,
			Extension: extension,
			CreatedOn: createdOn,
			Size:      size,
		},
		SavedLocation: savedLocation,
	}
}

func (r *Repository) queryForResourceInformation(appId string, resourceId string) *sql.Row {
	return r.db.QueryRow(`
		SELECT r.id, r.name, r.extension, r.size, r.created_on, rr.saved_location
		FROM resources r 
		JOIN resource_relations rr ON r.id = rr.resource_id
		WHERE rr.app_id = $1 AND r.id = $2
	`, appId, resourceId)
}
