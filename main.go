package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go-rest-api/db"
	"github.com/jackc/pgx/v4/pgxpool"
)

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World!")
    fmt.Println("Endpoint Hit: homePage")
}

func getUsers(w http.ResponseWriter, r *http.Request, dbpool *pgxpool.Pool) {
	users := db.GetUsers(dbpool)
	fmt.Fprintf(w, users)
}

func handleRequests(dbpool *pgxpool.Pool) {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		getUsers(w, r, dbpool)
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
