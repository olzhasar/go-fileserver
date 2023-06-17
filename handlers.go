package main

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const UPLOAD_URL = "/upload"
const DOWNLOAD_URL = "/download"

const MSG_UPLOAD_SUCCESS = "File uploaded successfully"
const MSG_ERR_INVALID_REQUEST_METHOD = "Invalid request method"
const MSG_ERR_CANNOT_READ_FILE = "Unable to read uploaded file"
const MSG_ERR_FILE_NOT_FOUND = "File not found"
const MSG_ERR_CANNOT_SEND_FILE = "Unable to send file"
const MSG_ERR_MISSING_QUERY_PARAM = "Missing filename query param"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to the FileServer. Use upload/ or download/ endpoints")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, MSG_ERR_INVALID_REQUEST_METHOD, http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_READ_FILE, http.StatusBadRequest)
		return
	}

	if err := saveFile(fileHeader.Filename, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, MSG_UPLOAD_SUCCESS)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, MSG_ERR_INVALID_REQUEST_METHOD, http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	fileName := queryParams.Get("filename")

	if fileName == "" {
		http.Error(w, MSG_ERR_MISSING_QUERY_PARAM, http.StatusBadRequest)
		return
	}

	filePath := getUploadFilePath(fileName)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, MSG_ERR_FILE_NOT_FOUND, http.StatusNotFound)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_READ_FILE, http.StatusInternalServerError)
		return
	}

	setFileHeaders(w, file, stat.Size(), fileName, filePath)

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_SEND_FILE, http.StatusInternalServerError)
		return
	}
}

func setFileHeaders(w http.ResponseWriter, file *os.File, size int64, name string, path string) {
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Disposition", "attachment; filename="+name)

	contentType := guessFileContentType(path)
	w.Header().Set("Content-Type", contentType)
}

func guessFileContentType(path string) string {
	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType
}
