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
	"time"
)

const PORT = "8080"
const UPLOAD_DIR = "uploads"

type LoggingMiddleware struct {
	handler http.Handler
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	l.handler.ServeHTTP(w, r)

	log.Printf("%s %s DURATION: %v", r.Method, r.URL, time.Since(start))
}

func MakeLoggedHandler(handler http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{handler}
}

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

	defer file.Close()

	newFilePath := getUploadFilePath(fileHeader.Filename)
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

func getUploadFilePath(filename string) string {
	return filepath.Join(UPLOAD_DIR, filename)
}
