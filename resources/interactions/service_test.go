package interactions

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mensurowary/juno/config"
	"github.com/mensurowary/juno/resources/download"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var tmpDir = "./tmp"

func TestService_DeleteSingleResourceById(t *testing.T) {
	t.Run("When resource does not exist", func(t *testing.T) {
		db, _ := getDbAndMock(t)
		s := getService(db, download.NoDownloadableResource)
		result := s.DeleteSingleResourceById("123456", "admin")
		assert.Equal(t, CouldNotFind, result)
		t.Cleanup(func() {
			_ = db.Close()
		})
	})

	t.Run("When resource exists and successfully deleted", func(t *testing.T) {
		baseFilename := createDummyFile(t)
		db, mock := getDbAndMock(t)
		s := getService(db, download.DownloadableResource{
			SavedLocation: baseFilename,
		})

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("admin", "123456789").
			WillReturnResult(sqlmock.NewResult(-1, 1))
		mock.ExpectCommit()

		result := s.DeleteSingleResourceById("123456789", "admin")

		assert.Nil(t, result)
		assert.Nil(t, mock.ExpectationsWereMet())

		t.Cleanup(func() {
			_ = os.RemoveAll("./tmp")
			_ = db.Close()
		})
	})

	t.Run("When resource exists but data could not be deleted", func(t *testing.T) {
		baseFilename := createDummyFile(t)
		db, mock := getDbAndMock(t)
		s := getService(db, download.DownloadableResource{
			SavedLocation: baseFilename,
		})

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("admin", "123456789").
			WillReturnResult(sqlmock.NewResult(-1, 2)) // more than one rows affected
		mock.ExpectRollback()

		result := s.DeleteSingleResourceById("123456789", "admin")

		assert.Equal(t, CouldNotDeleteData, result)
		assert.Nil(t, mock.ExpectationsWereMet())

		t.Cleanup(func() {
			_ = os.RemoveAll("./tmp")
			_ = db.Close()
		})
	})

	t.Run("When resource exists and data successfully deleted but file could not be deleted", func(t *testing.T) {
		baseFilename := createDummyFile(t)

		db, mock := getDbAndMock(t)
		s := getService(db, download.DownloadableResource{
			SavedLocation: baseFilename + "random",
		})

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("admin", "123456789").
			WillReturnResult(sqlmock.NewResult(-1, 1))
		mock.ExpectCommit()

		result := s.DeleteSingleResourceById("123456789", "admin")

		assert.Equal(t, CouldNotDeleteFile, result)
		assert.Nil(t, mock.ExpectationsWereMet())

		t.Cleanup(func() {
			_ = os.RemoveAll("./tmp")
			_ = db.Close()
		})
	})
}

func TestRepository_DeleteResourceById(t *testing.T) {
	t.Run("Tx init fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := NewRepository(db)

		mock.ExpectBegin().WillReturnError(errors.New("begin resulted in error"))
		err := r.DeleteResourceById("", "")

		assert.Equal(t, ErrCouldNotStartTx, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepared Statement init fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			WillReturnError(errors.New("failed for no reason :)"))
		err := r.DeleteResourceById("", "")

		assert.Equal(t, ErrCouldNotCreatePS, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Prepared Statement exec fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("app_id", "resource_id").
			WillReturnError(errors.New("failed for no reason too"))
		err := r.DeleteResourceById("resource_id", "app_id")

		assert.Equal(t, ErrCouldNotExecStmt, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Reading rows affected fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("app_id", "resource_id").
			WillReturnResult(sqlmock.NewErrorResult(errors.New("fails for no reason")))
		err := r.DeleteResourceById("resource_id", "app_id")

		assert.Equal(t, ErrCouldNotReadRowsAffected, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Committing fails", func(t *testing.T) {
		db, mock := getDbAndMock(t)
		r := NewRepository(db)

		mock.ExpectBegin()
		mock.ExpectPrepare(`^DELETE FROM resources WHERE id IN*`).
			ExpectExec().
			WithArgs("app_id", "resource_id").
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(errors.New("fails for no reason"))
		err := r.DeleteResourceById("resource_id", "app_id")

		assert.Equal(t, ErrCouldNotCommit, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}

func getDbAndMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	return db, mock
}

func getService(db *sql.DB, expected download.DownloadableResource) *Service {
	rs := mockResourceService{expected}
	r := NewRepository(db)
	return NewService(r, &rs)
}

func createDummyFile(t *testing.T) string {
	config.Config.FileUploadDir = tmpDir
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		if err := os.Mkdir(tmpDir, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}
	handle, err := os.OpenFile(tmpDir+"/hello.txt", os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Close()
	return filepath.Base(handle.Name())
}

type mockResourceService struct {
	resource download.DownloadableResource
}

func (m *mockResourceService) GetSingleResourceInformation(params download.SingleResourceRequestParams) download.DownloadableResource {
	return m.resource
}
