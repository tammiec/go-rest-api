package model

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestGetDb(t *testing.T) {
	url, _ := os.LookupEnv("DATABASE_URL")
	db := GetDb(url)
	require.IsType(t, &sql.DB{}, db)
}

func TestGetUsersSuccessfully(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email", "password"})
	rows.AddRow("1", "Kaladin", "k@s.com", "password")
	rows.AddRow("2", "Adolin", "a@k.com", "password")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := GetUsers(db)

	require.NoError(t, err)
	require.Equal(t, 1, result[0].Id)
	require.Equal(t, "Kaladin", result[0].Name)
	require.Equal(t, "k@s.com", result[0].Email)
	require.Equal(t, "password", result[0].Password)
	require.Equal(t, "Adolin", result[1].Name)
	require.Equal(t, "a@k.com", result[1].Email)
	require.Equal(t, "password", result[1].Password)
}

func TestGetFieldSettingsQueryError(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Mock Error"))
	_, err := GetUsers(db)
	
	require.Equal(t, "Mock Error", err.Error())
}

func TestGetFieldSettingsQueryBadRow(t *testing.T) {
	db, mock := getMockDB(t)
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email", "password"})
	rows.AddRow("1", "Kaladin", 1, "password")
	rows.AddRow("2", "Adolin", "a@k.com", nil)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err := GetUsers(db)

	require.Equal(t, "sql: Scan error on column index 3, name \"password\": converting NULL to string is unsupported", err.Error())
}
