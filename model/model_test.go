package model

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
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

func TestGetDb(t *testing.T) {
	url, _ := os.LookupEnv("DATABASE_URL")
	db := GetDb(url)
	require.IsType(t, &sql.DB{}, db)
}

func TestGetUsersSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	rows.AddRow("2", "Adolin", "a@k.com")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result, err := GetUsers(db)

	require.NoError(t, err)
	require.Equal(t, 1, result[0].Id)
	require.Equal(t, "Kaladin", result[0].Name)
	require.Equal(t, "k@s.com", result[0].Email)
	require.Equal(t, "Adolin", result[1].Name)
	require.Equal(t, "a@k.com", result[1].Email)
}

func TestGetUsersQueryError(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Mock Error"))

	_, err := GetUsers(db)

	require.Error(t, err)
	require.Equal(t, "Mock Error", err.Error())
}

func TestGetUsersQueryBadRow(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", nil)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	_, err := GetUsers(db)

	require.Error(t, err)
	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
}

func TestGetUserSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	result, err := GetUser(db, 1)

	require.NoError(t, err)
	require.Equal(t, 1, result.Id)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestGetUserQueryError(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("Mock Error"))

	_, err := GetUser(db, 1)

	require.Error(t, err)
	require.Equal(t, "Mock Error", err.Error())
}

func TestGetUserQueryBadRow(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", nil)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	_, err := GetUser(db, 1)

	require.Error(t, err)
	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
}

func TestDeleteUserSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnRows(rows)

	result, err := DeleteUser(db, 1)

	require.NoError(t, err)
	require.Equal(t, 1, result.Id)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestDeleteUserQueryInvalidUser(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("DELETE").WithArgs(2).WillReturnError(errors.New("sql: no rows in result set"))

	_, err := DeleteUser(db, 2)

	require.Error(t, err)
	require.Equal(t, "sql: no rows in result set", err.Error())
}

func TestCreateUserSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("INSERT")
	rows := mock.NewRows([]string{"name", "email", "password"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("INSERT").WithArgs("Kaladin", "k@s.com", "password").WillReturnRows(rows)

	result, err := CreateUser(db, "Kaladin", "k@s.com", "password")

	require.NoError(t, err)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestUpdateUserSuccessfully(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnRows(rows)

	result, err := UpdateUser(db, 1, "Kaladin", "k@s.com", "password")

	require.NoError(t, err)
	require.Equal(t, "Kaladin", result.Name)
	require.Equal(t, "k@s.com", result.Email)
}

func TestUpdateUserInvalidUser(t *testing.T) {
	db, mock := getMockDB()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 2).WillReturnError(errors.New("sql: no rows in result set"))

	_, err := UpdateUser(db, 2, "Kaladin", "k@s.com", "password")

	require.Error(t, err)
	require.Equal(t, "sql: no rows in result set", err.Error())
}
