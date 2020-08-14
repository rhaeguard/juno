package interactions

import (
	"database/sql"
	"github.com/mensurowary/juno/resources/download"
)

func NewService(db *sql.DB, downloadService download.Service) Service {
	return Service{
		r: Repository{
			db: db,
		},
		resourceService: downloadService,
	}
}
