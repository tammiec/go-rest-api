package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"

	"go-rest-api/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

func getUsers(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	users := db.GetUsers(dbpool)
	var body []byte
	body, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func getUser(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool, id int) {
	user := db.GetUser(dbpool, id)
	var body []byte
	body, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(body)
}

func handleRequests(dbpool *pgxpool.Pool) {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		getUsers(w, r, dbpool)
	})
	http.HandleFunc("/users/5", func(w http.ResponseWriter, r *http.Request) {
		getUser(w, r, dbpool, 5)
	})

	fmt.Println("Now listening at port 8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
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
	dbpool := db.GetDbPool(dbUrl)
	defer dbpool.Close()
	
	handleRequests(dbpool)
}
