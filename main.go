package main

import (
	"github.com/olzhasar/go-fileserver/router"
	"github.com/olzhasar/go-fileserver/storages"
	"log"
	"net/http"
)

const UPLOAD_DIR = "uploads"
const PORT = "8080"

func main() {
	storage := storages.NewFileSystemStoage(UPLOAD_DIR)
	router := router.NewRouter(storage)

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", router.UploadHandler)
	mux.HandleFunc("/download", router.DownloadHandler)
	mux.HandleFunc("/", router.RootHandler)

	loggedMux := MakeLoggedHandler(mux)

	log.Printf("Starting the server on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, loggedMux))
}
