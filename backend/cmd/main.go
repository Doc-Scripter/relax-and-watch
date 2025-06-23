package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Serve static files from the "frontend/public" directory
func main() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("../../frontend/public"))
	r.PathPrefix("/").Handler(fs)

	fmt.Println("Server starting on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
