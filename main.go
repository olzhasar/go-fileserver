package main

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const PORT = "8080"

const MSG_INVALID_REQUEST_METHOD = "Invalid request method"
const MSG_CANNOT_READ_FILE = "Unable to read uploaded file"

func checkUploadDir() {
	err := os.MkdirAll(UPLOAD_DIR, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
}

func main() {
	checkUploadDir()

	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/download", downloadHandler)
	mux.HandleFunc("/", rootHandler)

	loggedMux := MakeLoggedHandler(mux)

	log.Printf("Starting server on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, loggedMux))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to the FileServer. Use upload/ or download/ endpoints")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read uploaded file", http.StatusBadRequest)
		return
	}

	if err := saveFile(fileHeader.Filename, &file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File uploaded successfully")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	filename := queryParams.Get("filename")

	if filename == "" {
		http.Error(w, "Missing filename query parameter", http.StatusBadRequest)
		return
	}

	filePath := getUploadFilePath(filename)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Unable to read file info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(fileinfo.Size(), 10))

	contentType := mime.TypeByExtension(filepath.Ext(filePath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Unable to send file", http.StatusInternalServerError)
		return
	}
}
