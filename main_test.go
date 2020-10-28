package main

import (
	"database/sql"
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

func TestGetEnvOk(t *testing.T) {
	varname := "test"
	os.Setenv(varname, "1")
	v := getEnv(varname)
	require.Equal(t, "1", v)
}

func TestGetEnvEmptyString(t *testing.T) {
	varname := "test"
	os.Setenv(varname, "")
	require.Panics(t, func() { getEnv(varname) })
}

func TestGetEnvNotSet(t *testing.T) {
	varname := "test"
	os.Unsetenv(varname)
	require.Panics(t, func() { getEnv(varname) })
}

// func httpRequest(method string, url string, headers map[string]string) ([]byte, *http.Response, error) {
// 	request := httptest.NewRequest(method, url, nil)
// 	for key, val := range headers {
// 		request.Header.Set(key, val)
// 	}
// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, request)
// 	resp := recorder.Result()
// 	body, err := ioutil.ReadAll(recorder.Body)
// 	return body, resp, err
// }
