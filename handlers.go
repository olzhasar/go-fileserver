package fileserver

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
	filename := queryParams.Get("filename")

	if filename == "" {
		http.Error(w, MSG_ERR_MISSING_QUERY_PARAM, http.StatusBadRequest)
		return
	}

	filePath := getUploadFilePath(filename)

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, MSG_ERR_FILE_NOT_FOUND, http.StatusNotFound)
		return
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_READ_FILE, http.StatusInternalServerError)
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
		http.Error(w, MSG_ERR_CANNOT_SEND_FILE, http.StatusInternalServerError)
		return
	}
}
