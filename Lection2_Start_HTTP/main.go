package main

import (
	"fmt"
	"log"
	"net/http"
)

func GetGreet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Heya!</h1>")
	fmt.Println("Mtd: ", r.Method)
	fmt.Println("URL:", r.URL)
}

func main() {
	http.HandleFunc("GET /", GetGreet)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
