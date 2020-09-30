package model

import (
	"database/sql"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "log"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/jackc/pgx/v4/pgxpool"
)

func getMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	return db, mock
}

func TestGetDbPool(t *testing.T) {
	v, _ := os.LookupEnv("DATABASE_URL")
	db := GetDbPool(v)
	require.IsType(t, &pgxpool.Pool{}, db)
}

func TestGetUsers(t *testing.T) {
	db, mock := getMockDB(t)

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email", "password"})
	rows.AddRow("1", "Kaladin", "k@s.com", "password")
	rows.AddRow("2", "Adolin", "a@k.com", "password")
	rows.AddRow("3", "Shallan", "s@d.com", "password")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	result := GetUsers(db)
}
