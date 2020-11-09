package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tammiec/go-rest-api/model"
)

var (
	httpReadTimeout  = 5 * time.Second
	httpWriteTimeout = 5 * time.Second
	httpIdleTimeout  = 1 * time.Minute
)

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: add a DB ping check
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hello World")
}

func getUsersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	users, err := model.GetUsers(db)
	if err != nil {
		log.Println(err)
		// throw custom error type since psql doesn't throw a sql.ErrNoRows when there is an empty array
		// TODO - see if there's a better way to handle this and keep it consistent with the switch/case
		if err.Error() == "no users found" {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}
	marshalAndWriteJson(users, w)
}

func getUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, id int) {
	user, err := model.GetUser(db, id)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			http.Error(w, err.Error(), 404)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}
	marshalAndWriteJson(user, w)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, id int) {
	user, err := model.DeleteUser(db, id)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			http.Error(w, err.Error(), 404)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}
	marshalAndWriteJson(user, w)
}

func createUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, name string, email string, password string) {
	user, err := model.CreateUser(db, name, email, password)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			http.Error(w, err.Error(), 404)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}
	marshalAndWriteJson(user, w)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, id int, name string, email string, password string) {
	user, err := model.UpdateUser(db, id, name, email, password)
	if err != nil {
		log.Println(err)
		switch err {
		case sql.ErrNoRows:
			http.Error(w, err.Error(), 404)
		default:
			http.Error(w, err.Error(), 500)
		}
		return
	}
	marshalAndWriteJson(user, w)
}

func marshalAndWriteJson(data interface{}, w http.ResponseWriter) {
	var body []byte
	body, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func validateId(idString string, w http.ResponseWriter) int {
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	}
	return id
}

func parseRequest(r *http.Request) (string, string, string) {
	r.ParseForm()
	name := r.Form["name"][0]
	email := r.Form["email"][0]
	password := r.Form["password"][0]
	return name, email, password
}

func getRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/readiness", readinessHandler).Methods(http.MethodGet)
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsersHandler(w, r, db)
		} else if r.Method == http.MethodPost {
			name, email, password := parseRequest(r)
			createUserHandler(w, r, db, name, email, password)
		}
	}).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := validateId(mux.Vars(r)["id"], w)
		if r.Method == http.MethodGet {
			getUserHandler(w, r, db, id)
		} else if r.Method == http.MethodDelete {
			deleteUserHandler(w, r, db, id)
		} else if r.Method == http.MethodPut {
			name, email, password := parseRequest(r)
			updateUserHandler(w, r, db, id, name, email, password)
		}
	}).Methods(http.MethodGet, http.MethodDelete, http.MethodPut)

	return router
}

func httpServer(host string, port string, db *sql.DB) {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		IdleTimeout:  httpIdleTimeout,
		Handler:      getRouter(db),
	}
	log.Printf("Listening http://%s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func getEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Variable %s is not set", key))
	} else if v == "" {
		panic(fmt.Sprintf("Variable %s is blank", key))
	}
	return v
}

func main() {
	dbUrl := getEnv("DATABASE_URL")
	httpHost := getEnv("HTTP_HOST")
	httpPort := getEnv("HTTP_PORT")

	db := model.GetDb(dbUrl)
	defer db.Close()

	httpServer(httpHost, httpPort, db)
}
