package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "Hello World!")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
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
	// user := getEnv("PSQL_USER")
	// password := getEnv("PSQL_PASSWORD")
	// host := getEnv("PSQL_HOSTNAME")
	// dbPort := getEnv("PSQL_PORT")
	// dbName := getEnv("PSQL_DB_NAME")
	// httpHost := getEnv("HTTP_HOST")
	// httpPort := getEnv("HTTP_PORT")
	
    handleRequests()
}
