package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tammiec/go-rest-api/model"
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
	dbpool := model.GetDbPool(getEnv('DATABASE_URL'))
	defer dbpool.Close()

	var name string
	err = dbpool.QueryRow(context.Background(), "SELECT name FROM users WHERE id=5").Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(name)
	
	handleRequests()
}
