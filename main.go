package main


import (
	"fmt"
	"log"
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

func main() {
    handleRequests()
}
