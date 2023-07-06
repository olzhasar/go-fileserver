package main

import (
	"log"
	"net/http"

	"github.com/olzhasar/go-fileserver/loggers"
	"github.com/olzhasar/go-fileserver/manager"
	"github.com/olzhasar/go-fileserver/middleware"
	"github.com/olzhasar/go-fileserver/registry"
	"github.com/olzhasar/go-fileserver/server"
	"github.com/olzhasar/go-fileserver/storages"
)

const UPLOAD_DIR = "uploads"
const PORT = "8080"

func main() {
	storage := storages.NewFileSystemStoage(UPLOAD_DIR)
	registry, err := registry.NewSQLiteRegistry("./db.sqlite3")
	if err != nil {
		log.Fatalf("Error while initializing SQLite registry\n%s", err)
	}

	mgr := manager.NewFileManager(registry, storage)
	server := server.NewFileServer(mgr)

	logger := &loggers.StdLogger{}
	loggedServer := middleware.MakeLoggedHandler(server, logger)

	log.Printf("Starting the server on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, loggedServer))
}
