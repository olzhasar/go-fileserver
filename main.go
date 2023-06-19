package main

import (
	"log"
	"net/http"
)

var UPLOAD_DIR = "uploads"
var storage Storage

const PORT = "8080"

func main() {
	storage = NewFileSystemStoage(UPLOAD_DIR)

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/", rootHandler)

	loggedMux := MakeLoggedHandler(mux)

	log.Printf("Starting the server on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, loggedMux))
}
