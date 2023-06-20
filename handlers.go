package main

import (
	"fmt"
	"io"
	"net/http"
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

	if err := storage.saveFile(fileHeader.Filename, file); err != nil {
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

	upload, err := storage.loadFile(fileName)
	if err != nil {
		http.Error(w, MSG_ERR_FILE_NOT_FOUND, http.StatusNotFound)
		return
	}
	defer upload.file.Close()

	_, err = io.Copy(w, upload.file)
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_SEND_FILE, http.StatusInternalServerError)
		return
	}

	setFileHeaders(w, upload)
}

func setFileHeaders(w http.ResponseWriter, upload UploadedFile) {
	w.Header().Set("Content-Length", strconv.FormatInt(upload.size, 10))
	w.Header().Set("Content-Disposition", "attachment; filename="+upload.name)
	w.Header().Set("Content-Type", guessFileContentType(upload))
}

func guessFileContentType(upload UploadedFile) string {
	contentType := upload.MimeTypeByExt()
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType
}
