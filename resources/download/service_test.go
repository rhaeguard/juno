package download

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mensurowary/juno/commons"
	"github.com/mensurowary/juno/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path/filepath"
	"testing"
	"time"
)

func TestService_GetAppResourcesInformation(t *testing.T) {
	appId := "admin"

	columns := []string{
		"id", "name", "extension", "size", "created_on",
	}

	t.Run("Successfully retrieves the app resources data by app id", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		expected1 := Resource{
			Id:        "123456",
			Name:      "Mipe Stiopic",
			Extension: "txt",
			Size:      123456,
			CreatedOn: time.Now(),
		}

		expected2 := Resource{
			Id:        "654321",
			Name:      "Daniel Cormier",
			Extension: "epub",
			Size:      654321,
			CreatedOn: time.Date(2020, time.March, 14, 12, 6, 0, 0, time.Local),
		}

		mock.ExpectQuery(`^SELECT (.+) FROM resources*`).
			WithArgs("admin").
			WillReturnRows(
				sqlmock.
					NewRows(columns).
					AddRow(spread(expected1)...).
					AddRow(spread(expected2)...),
			)

		s := getService(db)

		ri := s.GetAppResourcesInformation(appId)

		assert.Nil(t, ri.Err)
		assert.NotNil(t, ri.Resources)
		assert.Equal(t, []Resource{expected1, expected2}, ri.Resources)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Response should be empty slice when there is no data", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		mock.ExpectQuery(`^SELECT (.+) FROM resources*`).
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows(columns))

		s := getService(db)

		ri := s.GetAppResourcesInformation(appId)

		assert.Nil(t, ri.Err)
		assert.NotNil(t, ri.Resources)
		assert.Empty(t, ri.Resources)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Response should be empty slice and error should exist when query value mapping fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		mock.ExpectQuery(`^SELECT (.+) FROM resources*`).
			WithArgs("admin").
			WillReturnRows(sqlmock.NewRows(columns).AddRow("123654", "name", "ext", nil, time.Now()).RowError(1, errors.New("could not do stuff")))

		s := getService(db)

		ri := s.GetAppResourcesInformation(appId)

		assert.Equal(t, CouldNotRetrieveResults, ri.Err)
		assert.NotNil(t, ri.Resources)
		assert.Empty(t, ri.Resources)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Response should be empty slice and error should exist when query initialization", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		mock.ExpectQuery(`^SELECT (.+) FROM resources*`).
			WithArgs("admin").
			WillReturnError(errors.New("query is incorrect or whatever"))

		s := getService(db)

		ri := s.GetAppResourcesInformation(appId)

		assert.Equal(t, CouldNotRetrieveResults, ri.Err)
		assert.NotNil(t, ri.Resources)
		assert.Empty(t, ri.Resources)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestService_GetSingleResourceInformation(t *testing.T) {
	columns := []string{
		"id", "name", "extension", "size", "created_on", "saved_location",
	}
	t.Run("Successfully gets single resource information", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		expected := DownloadableResource{
			Resource: Resource{
				Id:        "123456",
				Name:      "mock",
				Extension: "pdf",
				CreatedOn: time.Now(),
				Size:      123456,
			},
			SavedLocation: "./hello/mock.pdf",
		}

		mock.ExpectQuery(`\s*SELECT (.+) \s*FROM resources r*`).
			WithArgs("admin", "123456").
			WillReturnRows(sqlmock.NewRows(columns).AddRow(
				spreadDR(expected)...,
			))

		s := getService(db)

		actual := s.GetSingleResourceInformation(SingleResourceRequestParams{
			ResourceId: "123456",
			AppId:      "admin",
		})

		assert.Equal(t, expected, actual)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty result when no data is available", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		defer db.Close()

		mock.ExpectQuery(`\s*SELECT (.+) \s*FROM resources r*`).
			WithArgs("admin", "123456").
			WillReturnRows(sqlmock.NewRows(columns))

		s := getService(db)

		actual := s.GetSingleResourceInformation(SingleResourceRequestParams{
			ResourceId: "123456",
			AppId:      "admin",
		})

		assert.Equal(t, NoDownloadableResource, actual)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestService_GetSingleResource(t *testing.T) {
	columns := []string{
		"id", "name", "extension", "size", "created_on", "saved_location",
	}

	t.Run("Get single resource information", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		s := getService(db)

		dr := DownloadableResource{
			Resource: Resource{
				Id:        "123456789",
				Name:      "mock",
				Extension: "pdf",
				CreatedOn: time.Now(),
				Size:      123456,
			},
			SavedLocation: "./hello/mock.pdf",
		}
		mock.ExpectQuery(`\s*SELECT (.+) \s*FROM resources r*`).
			WithArgs("admin", "123456789").
			WillReturnRows(sqlmock.NewRows(columns).AddRow(
				spreadDR(dr)...,
			))

		result := s.GetSingleResource(SingleResourceRequestParams{
			AppId:      "admin",
			ResourceId: "123456789",
			Download:   false,
		})

		assert.Equal(t, SingleResourceResult{
			File:   nil,
			Data:   commons.MakeSuccessResponse("Successfully retrieved the resource information", dr.Resource),
			Status: http.StatusOK,
		}, result)
	})

	t.Run("Get single resource download information", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		s := getService(db)

		dr := DownloadableResource{
			Resource: Resource{
				Id:        "123456789",
				Name:      "mock",
				Extension: "pdf",
				CreatedOn: time.Now(),
				Size:      123456,
			},
			SavedLocation: "hello/mock.pdf",
		}
		mock.ExpectQuery(`\s*SELECT (.+) \s*FROM resources r*`).
			WithArgs("admin", "123456789").
			WillReturnRows(sqlmock.NewRows(columns).AddRow(
				spreadDR(dr)...,
			))

		result := s.GetSingleResource(SingleResourceRequestParams{
			AppId:      "admin",
			ResourceId: "123456789",
			Download:   true,
		})

		assert.Equal(t, SingleResourceResult{
			File: &SingleResourceFileResult{
				Name: "mock.pdf",
				Path: filepath.Join(config.Config.FileUploadDir, "hello/mock.pdf"),
			},
		}, result)
	})

	t.Run("Return empty results when no data is available", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		s := getService(db)

		mock.ExpectQuery(`\s*SELECT (.+) \s*FROM resources r*`).
			WithArgs("admin", "123456789").
			WillReturnRows(sqlmock.NewRows(columns))

		result := s.GetSingleResource(SingleResourceRequestParams{
			AppId:      "admin",
			ResourceId: "123456789",
			Download:   true,
		})

		assert.Equal(t, SingleResourceResult{
			File:   nil,
			Data:   commons.MakeFailureResponse("Could not find the requested resource", http.StatusNotFound),
			Status: http.StatusNotFound,
		}, result)
	})
}

func TestGetFileName(t *testing.T) {
	tt := []struct {
		RequestedName, SavedName, Extension string
		Expected                            string
	}{
		{
			RequestedName: "hello",
			SavedName:     "hey",
			Extension:     "hi",
			Expected:      "hello.hi",
		},
		{
			RequestedName: "hello",
			SavedName:     "hey",
			Extension:     "",
			Expected:      "hello",
		},
	}

	for _, tc := range tt {
		params := SingleResourceRequestParams{Name: tc.RequestedName}

		resource := Resource{
			Name:      tc.SavedName,
			Extension: tc.Extension,
		}

		result := getFileName(&params, &resource)
		assert.Equal(t, tc.Expected, result)
	}

}

func getDbAndMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	return db, mock
}

func getService(db *sql.DB) *Service {
	r := NewRepository(db)
	return NewService(r)
}

func spread(resource Resource) []driver.Value {
	return []driver.Value{
		resource.Id, resource.Name, resource.Extension, resource.Size, resource.CreatedOn,
	}
}

func spreadDR(dr DownloadableResource) []driver.Value {
	values := spread(dr.Resource)
	return append(values, dr.SavedLocation)
}
