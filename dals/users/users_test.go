package users

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tammiec/go-rest-api/testutils"
)

var db *sql.DB

func TestMain(m *testing.M) {
	db = testutils.GetDB()
	defer db.Close()
	os.Exit(m.Run())
}

func usersSetup() {
	testutils.RunSQLFile(db, `users_schema.sql`)
}

// func getMockDB() (*sql.DB, sqlmock.Sqlmock) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
// 	}
// 	return db, mock
// }

// // func TestGetDb(t *testing.T) {
// // 	url, _ := os.LookupEnv("DATABASE_URL")
// // 	db := GetDb(url)
// // 	require.IsType(t, &sql.DB{}, db)
// // }

func TestUsers_ListSuccess(t *testing.T) {
	usersSetup()
	testutils.RunSQLFile(db, `users_insert.sql`)

	usersDAL := impl{db: db}

	users, err := usersDAL.List()

	require.NoError(t, err)
	require.Equal(t, "", users)
}

// func TestListQueryError(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectQuery("SELECT").WillReturnError(errors.New("Mock Error"))

// 	_, err := List(db)

// 	require.Error(t, err)
// 	require.Equal(t, "Mock Error", err.Error())
// }

// func TestListQueryBadRow(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow("1", "Kaladin", nil)
// 	mock.ExpectQuery("SELECT").WillReturnRows(rows)

// 	_, err := List(db)

// 	require.Error(t, err)
// 	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
// }

// func TestGetSuccessfully(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("SELECT")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow("1", "Kaladin", "k@s.com")
// 	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

// 	result, err := Get(db, 1)

// 	require.NoError(t, err)
// 	require.Equal(t, 1, result.Id)
// 	require.Equal(t, "Kaladin", result.Name)
// 	require.Equal(t, "k@s.com", result.Email)
// }

// func TestGetQueryError(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("SELECT")
// 	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(errors.New("Mock Error"))

// 	_, err := Get(db, 1)

// 	require.Error(t, err)
// 	require.Equal(t, "Mock Error", err.Error())
// }

// func TestGetQueryBadRow(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("SELECT")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow("1", "Kaladin", nil)
// 	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

// 	_, err := Get(db, 1)

// 	require.Error(t, err)
// 	require.Equal(t, "sql: Scan error on column index 2, name \"email\": converting NULL to string is unsupported", err.Error())
// }

// func TestDeleteSuccessfully(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("DELETE")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow("1", "Kaladin", "k@s.com")
// 	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnRows(rows)

// 	result, err := Delete(db, 1)

// 	require.NoError(t, err)
// 	require.Equal(t, 1, result.Id)
// 	require.Equal(t, "Kaladin", result.Name)
// 	require.Equal(t, "k@s.com", result.Email)
// }

// func TestDeleteQueryInvalidUser(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("DELETE")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow("1", "Kaladin", "k@s.com")
// 	mock.ExpectQuery("DELETE").WithArgs(2).WillReturnError(errors.New("sql: no rows in result set"))

// 	_, err := Delete(db, 2)

// 	require.Error(t, err)
// 	require.Equal(t, "sql: no rows in result set", err.Error())
// }

// func TestCreateSuccessfully(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("INSERT")
// 	rows := mock.NewRows([]string{"name", "email", "password"})
// 	rows.AddRow(1, "Kaladin", "k@s.com")
// 	mock.ExpectQuery("INSERT").WithArgs("Kaladin", "k@s.com", "password").WillReturnRows(rows)

// 	result, err := Create(db, "Kaladin", "k@s.com", "password")

// 	require.NoError(t, err)
// 	require.Equal(t, "Kaladin", result.Name)
// 	require.Equal(t, "k@s.com", result.Email)
// }

// func TestUpdateSuccessfully(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("UPDATE")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow(1, "Kaladin", "k@s.com")
// 	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnRows(rows)

// 	result, err := Update(db, 1, "Kaladin", "k@s.com", "password")

// 	require.NoError(t, err)
// 	require.Equal(t, "Kaladin", result.Name)
// 	require.Equal(t, "k@s.com", result.Email)
// }

// func TestUpdateInvalidUser(t *testing.T) {
// 	db, mock := getMockDB()
// 	defer db.Close()

// 	mock.ExpectPrepare("UPDATE")
// 	rows := mock.NewRows([]string{"id", "name", "email"})
// 	rows.AddRow(1, "Kaladin", "k@s.com")
// 	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 2).WillReturnError(errors.New("sql: no rows in result set"))

// 	_, err := Update(db, 2, "Kaladin", "k@s.com", "password")

// 	require.Error(t, err)
// 	require.Equal(t, "sql: no rows in result set", err.Error())
// }
