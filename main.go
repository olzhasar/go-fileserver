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

const port = "8080"
const UPLOAD_DIR = "uploads"

func checkUploadDir() {
	err := os.MkdirAll(UPLOAD_DIR, 0755)
	if err != nil {
		log.Fatalf("Error creating directory: %v", err)
	}
}

func main() {
	checkUploadDir()

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/", rootHandler)

	fmt.Printf("Starting server on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
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

	defer file.Close()

	newFilePath := UPLOAD_DIR + "/" + fileHeader.Filename
	newFile, err := os.Create(newFilePath)

	if err != nil {
		http.Error(w, "Unable to create destination file", http.StatusInternalServerError)
		return
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Unable to save uploaded file", http.StatusInternalServerError)
		return
	}

	filename := fileHeader.Filename
	fmt.Printf("Filename %s\n", filename)

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

	filePath := UPLOAD_DIR + "/" + filename

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
