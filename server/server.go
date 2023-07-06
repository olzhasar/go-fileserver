package server

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/olzhasar/go-fileserver/manager"
	"github.com/olzhasar/go-fileserver/storages"
)

const UPLOAD_URL = "/upload"
const DOWNLOAD_URL = "/download"

const MSG_UPLOAD_SUCCESS = "File uploaded successfully"
const MSG_ERR_INVALID_REQUEST_METHOD = "Invalid request method"
const MSG_ERR_CANNOT_READ_FILE = "Unable to read uploaded file"
const MSG_ERR_FILE_NOT_FOUND = "File not found"
const MSG_ERR_CANNOT_SEND_FILE = "Unable to send file"
const MSG_ERR_MISSING_QUERY_PARAM = "Missing filename query param"

type FileServer struct {
	manager manager.SaverLoader
}

func NewFileServer(f manager.SaverLoader) *FileServer {
	return &FileServer{f}
}

func (f *FileServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", f.handleUpload)
	mux.HandleFunc("/download", f.handleDownload)
	mux.HandleFunc("/", f.handleRoot)

	mux.ServeHTTP(w, req)
}

func (f *FileServer) handleRoot(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Welcome to the FileServer. Use upload/ or download/ endpoints")
}

func (f *FileServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, MSG_ERR_INVALID_REQUEST_METHOD, http.StatusMethodNotAllowed)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_READ_FILE, http.StatusBadRequest)
		return
	}

	token, err := f.manager.SaveFile(fileHeader.Filename, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

	downloadUrl := buildDownloadURL(r.Host, token)
	fmt.Fprint(w, downloadUrl)
}

func (f *FileServer) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, MSG_ERR_INVALID_REQUEST_METHOD, http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()
	token := queryParams.Get("token")

	if token == "" {
		http.Error(w, MSG_ERR_MISSING_QUERY_PARAM, http.StatusBadRequest)
		return
	}

	upload, err := f.manager.LoadFile(token)
	if err != nil {
		http.Error(w, MSG_ERR_FILE_NOT_FOUND, http.StatusNotFound)
		return
	}
	defer upload.File.Close()

	_, err = io.Copy(w, upload.File)
	if err != nil {
		http.Error(w, MSG_ERR_CANNOT_SEND_FILE, http.StatusInternalServerError)
		return
	}

	setFileHeaders(w, upload)
}

func buildDownloadURL(host string, token string) string {
	return host + "?token=" + token
}

func setFileHeaders(w http.ResponseWriter, upload storages.UploadedFile) {
	w.Header().Set("Content-Length", strconv.FormatInt(upload.Size, 10))
	w.Header().Set("Content-Disposition", "attachment; filename="+upload.Name)
	w.Header().Set("Content-Type", guessFileContentType(upload))
}

func guessFileContentType(upload storages.UploadedFile) string {
	contentType := upload.MimeTypeByExt()
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return contentType
}
