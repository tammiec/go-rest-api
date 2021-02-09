package users

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func getMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	return db, mock
}

func TestListQueryError(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Mock Error"))

	mockDAL := &impl{db: db}

	_, err := mockDAL.List()

	require.Error(t, err)
	require.Equal(t, "Mock Error", err.Error())
}

func TestListQueryBadRow(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", nil)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	mockDAL := &impl{db: db}

	_, err := mockDAL.List()

	require.Error(t, err)
	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
}

func TestGetSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	mockDAL := &impl{db: db}

	result, err := mockDAL.Get(1)

	require.NoError(t, err)
	require.Equal(t, 1, result.Id)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestGetQueryError(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("Mock Error"))

	mockDAL := &impl{db: db}

	_, err := mockDAL.Get(1)

	require.Error(t, err)
	require.Equal(t, "Mock Error", err.Error())
}

func TestGetQueryBadRow(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", nil)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	mockDAL := &impl{db: db}

	_, err := mockDAL.Get(1)

	require.Error(t, err)
	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
}

func TestDeleteSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnRows(rows)

	mockDAL := &impl{db: db}

	result, err := mockDAL.Delete(1)

	require.NoError(t, err)
	require.Equal(t, 1, result.Id)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestDeleteQueryInvalidUser(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("DELETE").WithArgs(2).WillReturnError(errors.New("sql: no rows in result set"))

	mockDAL := &impl{db: db}

	_, err := mockDAL.Delete(2)

	require.Error(t, err)
	require.Equal(t, "sql: no rows in result set", err.Error())
}

func TestCreateSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("INSERT")
	rows := mock.NewRows([]string{"name", "email", "password"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("INSERT").WithArgs("Kaladin", "k@s.com", "password").WillReturnRows(rows)

	mockDAL := &impl{db: db}

	result, err := mockDAL.Create("Kaladin", "k@s.com", "password")

	require.NoError(t, err)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestUpdateSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnRows(rows)

	mockDAL := &impl{db: db}

	result, err := mockDAL.Update(1, "Kaladin", "k@s.com", "password")

	require.NoError(t, err)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestUpdateInvalidUser(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 2).WillReturnError(errors.New("sql: no rows in result set"))

	mockDAL := &impl{db: db}

	_, err := mockDAL.Update(2, "Kaladin", "k@s.com", "password")

	require.Error(t, err)
	require.Equal(t, "sql: no rows in result set", err.Error())
}
