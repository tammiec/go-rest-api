package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"strconv"
	"database/sql"

	"go-rest-api/model"
	"github.com/gorilla/mux"
)

func getUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	users, err := model.GetUsers(db)
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
	var body []byte
	body, err = json.Marshal(users)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func getUser(w http.ResponseWriter, r *http.Request, db *sql.DB, id int) {
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
	var body []byte
	body, err = json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func deleteUser(w http.ResponseWriter, r *http.Request, db *sql.DB, id int) {
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
	var body []byte
	body, err = json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func createUser(w http.ResponseWriter, r *http.Request, db *sql.DB, name string, email string, password string) {
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
	var body []byte
	body, err = json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func updateUser(w http.ResponseWriter, r *http.Request, db *sql.DB, id int, name string, email string, password string) {
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
	var body []byte
	body, err = json.Marshal(user)
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

func handleRequests(db *sql.DB) {
	router := mux.NewRouter()
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsers(w, r, db)
		} else if r.Method == http.MethodPost {
			name, email, password := parseRequest(r)
			createUser(w, r, db, name, email, password)
		}
	}).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/users/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := validateId(mux.Vars(r)["id"], w)
		if r.Method == http.MethodGet {
			getUser(w, r, db, id)
		} else if r.Method == http.MethodDelete {
			deleteUser(w, r, db, id)
		} else if r.Method == http.MethodPut {
			name, email, password := parseRequest(r)
			updateUser(w, r, db, id, name, email, password)
		}
	}).Methods(http.MethodGet, http.MethodDelete, http.MethodPut)

	fmt.Println("Now listening at port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
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
	db := model.GetDb(dbUrl)
	defer db.Close()
	
	handleRequests(db)
}
