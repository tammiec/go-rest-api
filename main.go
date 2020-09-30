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
	// "github.com/jackc/pgx/v4/pgxpool"
	"github.com/gorilla/mux"
)

// func getUsers(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
// 	users := db.GetUsers(dbpool)
// 	var body []byte
// 	body, err := json.Marshal(users)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Write(body)
// }

func getUser(w http.ResponseWriter, r *http.Request, db *sql.DB, id int) {
	user := model.GetUser(db, id)
	var body []byte
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

// func deleteUser(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool, id int) {
// 	user := db.DeleteUser(dbpool, id)
// 	var body []byte
// 	body, err := json.Marshal(user)
// 	if err != nil {
// 		log.Println(err)
// 		http.Error(w, err.Error(), 500)
// 		return
// 	}
// 	w.Write(body)
// }

func validateId(idString string, w http.ResponseWriter) int {
	id, err := strconv.Atoi(idString)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	}
	return id
}

func handleRequests(db *sql.DB) {
	router := mux.NewRouter()
	// router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
	// 	getUsers(w, r, dbpool)
	// })
	router.HandleFunc("/users/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		id := validateId(mux.Vars(r)["id"], w)
		if r.Method == http.MethodGet {
			getUser(w, r, db, id)
		} else if r.Method == http.MethodDelete {
			// deleteUser(w, r, dbpool, id)
		}
	}).Methods(http.MethodGet, http.MethodDelete)

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
