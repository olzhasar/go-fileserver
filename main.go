package main

import (
	"log"
	"net/http"
)

var UPLOAD_DIR = "uploads"

const PORT = "8080"

func main() {
	createDirIfNotExists(UPLOAD_DIR)

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/", rootHandler)

	loggedMux := MakeLoggedHandler(mux)

	log.Printf("Starting server on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, loggedMux))
}
