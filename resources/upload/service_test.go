package upload

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func TestService_HandleUpload(t *testing.T) {
	t.Run("Error occurred while uploading the file", func(t *testing.T) {
		db, _ := getDbAndMock(t)
		s := NewService(db)

		writer := &mockFileWriter{
			err: errors.New("failure"),
		}
		header := &multipart.FileHeader{}
		appId := "app_id"
		values := map[string][]string{
			"name": {"hello"},
		}

		resourceID, err := s.HandleUpload(writer, header, appId, values)

		assert.Equal(t, EmptyID, resourceID)
		assert.Equal(t, ErrFileCouldNotBeUploaded, err)
	})

	t.Run("Upload successful", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		s := NewService(db)

		mock.ExpectBegin()
		mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
			ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
			WillReturnResult(sqlmock.NewResult(-1, 1))

		mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
			ExpectExec().WithArgs("app_id", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(-1, 1))

		mock.ExpectCommit()

		writer := &mockFileWriter{}
		header := &multipart.FileHeader{
			Size:     123456,
			Filename: "hello.pdf",
		}
		appId := "app_id"
		values := map[string][]string{}

		resourceID, err := s.HandleUpload(writer, header, appId, values)

		assert.NotEmpty(t, resourceID)
		assert.Nil(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("When saving upload information fails return error", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		s := NewService(db)

		mock.ExpectBegin()
		mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
			WillReturnError(errors.New("a random error"))

		writer := &mockFileWriter{}
		header := &multipart.FileHeader{
			Size:     123456,
			Filename: "hello.pdf",
		}
		appId := "app_id"
		values := map[string][]string{}

		resourceID, err := s.HandleUpload(writer, header, appId, values)

		assert.Empty(t, resourceID)
		assert.Equal(t, ErrFileCouldNotBeUploaded, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func TestRepository_saveUploadedResourceInformation(t *testing.T) {
	t.Run("Successfully persists the uploaded resource data", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := &Repository{db}

		mock.ExpectBegin()
		mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
			ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
			WillReturnResult(sqlmock.NewResult(-1, 1))

		mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
			ExpectExec().WithArgs("admin", sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(-1, 1))

		mock.ExpectCommit()

		params := &SaveUploadedResourceParameters{
			FileName:          "hello",
			FileExtension:     "pdf",
			FileSize:          123456,
			AppID:             "admin",
			UploadDestination: "src/hello.pdf",
		}

		result := r.saveUploadedResourceInformation(params)

		assert.NotEmpty(t, result.ID)
		assert.Nil(t, result.Err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Committing fails when persisting the uploaded resource data", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
				ExpectExec().WithArgs("admin", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectCommit().WillReturnError(errors.New("commit failure due to ninja turtles"))
		})
	})

	t.Run("More than one row of uploaded resource info saved should rollback", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 2))

			mock.ExpectRollback()
		})
	})

	t.Run("More than one row of uploaded resource info saved should rollback and rollback fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 2))

			mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
		})
	})

	t.Run("More than one row of uploaded resource relations info saved should rollback", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
				ExpectExec().WithArgs("admin", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(-1, 2))

			mock.ExpectRollback()
		})
	})

	t.Run("More than one row of uploaded resource relations info saved should rollback and rollback fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
				ExpectExec().WithArgs("admin", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(-1, 2))

			mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))
		})
	})

	t.Run("Executing the resource relations insert fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
				ExpectExec().WillReturnError(errors.New("exec could not be performed"))
		})
	})

	t.Run("Preparing the resource relations insert statement fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WithArgs(sqlmock.AnyArg(), "hello", "pdf", 123456).
				WillReturnResult(sqlmock.NewResult(-1, 1))

			mock.ExpectPrepare("^INSERT INTO resource_relations(app_id, resource_id, saved_location)*").
				WillReturnError(errors.New("prep init failed"))
		})
	})

	t.Run("Executing the upload resource insert fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				ExpectExec().WillReturnError(errors.New("exec could not be performed"))
		})
	})

	t.Run("Preparing the upload resource info insert statement fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectPrepare("^INSERT INTO resources(id, name, extension, size, created_on)*").
				WillReturnError(errors.New("prep init failed"))
		})
	})

	t.Run("Transaction begin fails", func(t *testing.T) {
		runAndExpect(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin().WillReturnError(errors.New("tx init failed"))
		})
	})
}

type mockFileWriter struct {
	err error
}

func (m *mockFileWriter) SaveFileTo(file *multipart.FileHeader, dst string) error {
	return m.err
}

func getDbAndMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	return db, mock
}

func runAndExpect(t *testing.T, mockExpectations func(mock sqlmock.Sqlmock)) {
	db, mock := getDbAndMock(t)
	r := &Repository{db}

	mockExpectations(mock)

	params := &SaveUploadedResourceParameters{
		FileName:          "hello",
		FileExtension:     "pdf",
		FileSize:          123456,
		AppID:             "admin",
		UploadDestination: "src/hello.pdf",
	}

	result := r.saveUploadedResourceInformation(params)

	assert.Empty(t, result.ID)
	assert.Equal(t, errCouldNotPersist, result.Err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
