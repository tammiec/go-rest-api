package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/tammiec/go-rest-api/model"
)

func httpRequest(router *mux.Router, method string, url string, headers map[string]string) ([]byte, *http.Response, error) {
	request := httptest.NewRequest(method, url, nil)
	for key, val := range headers {
		request.Header.Set(key, val)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	resp := recorder.Result()
	body, err := ioutil.ReadAll(recorder.Body)
	return body, resp, err
}

func getMockDBAndRouter() (*sql.DB, sqlmock.Sqlmock, *mux.Router) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err))
	}
	router := getRouter(db)
	return db, mock, router
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

func TestHandleReadinessOk(t *testing.T) {
	db, _, router := getMockDBAndRouter()
	defer db.Close()

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/readiness", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "Ready", string(body))
}

func TestHandleGetUsersHttpOk(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	rows.AddRow("2", "Adolin", "a@k.com")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "[{\"Id\":1,\"Name\":\"Kaladin\",\"Email\":\"k@s.com\"},{\"Id\":2,\"Name\":\"Adolin\",\"Email\":\"a@k.com\"}]", string(body))

	result := [2]model.User{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, string(body))
	require.Equal(t, "Kaladin", result[0].Name)
}

func TestHandleGetUsersNoRows(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	rows := mock.NewRows([]string{"id", "name", "email"})
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode, string(body))
}

func TestHandleGetUsersSqlError(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleGetUserHttpOk(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "{\"Id\":1,\"Name\":\"Kaladin\",\"Email\":\"k@s.com\"}", string(body))

	result := &model.User{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, string(body))
	require.Equal(t, "Kaladin", result.Name)
}

func TestHandleGetUserNotFound(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	rows := mock.NewRows([]string{"id", "name", "email"})
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode, string(body))
}

func TestHandleGetUserSqlError(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("SELECT")
	mock.ExpectQuery("SELECT").WillReturnError(errors.New("error"))

	body, resp, err := httpRequest(router, http.MethodGet, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleDeleteUserHttpOk(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow("1", "Kaladin", "k@s.com")
	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodDelete, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "{\"Id\":1,\"Name\":\"Kaladin\",\"Email\":\"k@s.com\"}", string(body))

	result := &model.User{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, string(body))
	require.Equal(t, "Kaladin", result.Name)
}

func TestHandleDeleteUserNotFound(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodDelete, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode, string(body))
}

func TestHandleDeleteUserSqlError(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("DELETE")
	mock.ExpectQuery("DELETE").WithArgs(1).WillReturnError(errors.New("error"))

	body, resp, err := httpRequest(router, http.MethodDelete, "http://localhost:1234/users/1", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleCreateUserHttpOk(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("INSERT")
	rows := mock.NewRows([]string{"name", "email", "password"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("INSERT").WithArgs("Kaladin", "k@s.com", "password").WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodPost, "http://localhost:1234/users?name=Kaladin&email=k@s.com&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "{\"Id\":1,\"Name\":\"Kaladin\",\"Email\":\"k@s.com\"}", string(body))

	result := &model.User{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, string(body))
	require.Equal(t, "Kaladin", result.Name)
}

func TestHandleCreateUserBadArgs(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("INSERT")
	rows := mock.NewRows([]string{"name", "email", "password"})
	rows.AddRow(1, "Kaladin", 1)
	mock.ExpectQuery("INSERT").WithArgs("Kaladin", 1, "password").WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodPost, "http://localhost:1234/users?name=Kaladin&email=1&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleCreateUserSqlError(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("INSERT")
	mock.ExpectQuery("INSERT").WithArgs("Kaladin", "k@s.com", "password").WillReturnError(errors.New("error"))

	body, resp, err := httpRequest(router, http.MethodPost, "http://localhost:1234/users?name=Kaladin&email=k@s.com&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleUpdateUserHttpOk(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodPut, "http://localhost:1234/users/1?name=Kaladin&email=k@s.com&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, "{\"Id\":1,\"Name\":\"Kaladin\",\"Email\":\"k@s.com\"}", string(body))

	result := &model.User{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err, string(body))
	require.Equal(t, "Kaladin", result.Name)
}

func TestHandleUpdateUserBadArgs(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	rows.AddRow(1, "Kaladin", "k@s.com")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", nil, 1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodPut, "http://localhost:1234/users/1?name=Kaladin&email=k@s.com&password=", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleUpdateUserSqlError(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnError(errors.New("error"))

	body, resp, err := httpRequest(router, http.MethodPut, "http://localhost:1234/users/1?name=Kaladin&email=k@s.com&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode, string(body))
}

func TestHandleUpdateUserNotFound(t *testing.T) {
	db, mock, router := getMockDBAndRouter()
	defer db.Close()

	mock.ExpectPrepare("UPDATE")
	rows := mock.NewRows([]string{"id", "name", "email"})
	mock.ExpectQuery("UPDATE").WithArgs("Kaladin", "k@s.com", "password", 1).WillReturnRows(rows)

	body, resp, err := httpRequest(router, http.MethodPut, "http://localhost:1234/users/1?name=Kaladin&email=k@s.com&password=password", nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode, string(body))
}
