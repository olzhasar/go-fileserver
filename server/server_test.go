package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/olzhasar/go-fileserver/storages"
)

type StubFile struct {
	fileName string
	content  string
}

type StubFileManager struct {
	data map[string]StubFile
}

func (s *StubFileManager) SaveFile(fileName string, content io.Reader) (token string, err error) {
	token = "token"
	buf := new(strings.Builder)
	io.Copy(buf, content)
	s.data[token] = StubFile{fileName, buf.String()}
	return token, nil
}

func (s *StubFileManager) LoadFile(token string) (upload storages.UploadedFile, err error) {
	loaded, ok := s.data[token]
	if !ok {
		return storages.UploadedFile{}, errors.New(fmt.Sprintf("token %q is missing", token))
	}

	buf := &bytes.Buffer{}
	buf.WriteString(loaded.content)

	f := storages.InMemoryFile{Buffer: buf}
	size := int64(len(loaded.content))

	return storages.UploadedFile{File: f, Name: loaded.fileName, Size: size}, nil
}

func NewStubFileManager() *StubFileManager {
	data := make(map[string]StubFile)
	return &StubFileManager{data}
}

func TestRoot(t *testing.T) {
	mgr := NewStubFileManager()
	server := NewFileServer(mgr)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertResponseStatus(t, response, http.StatusOK)
}

func TestUpload(t *testing.T) {
	mgr := NewStubFileManager()
	server := NewFileServer(mgr)

	t.Run("uploads successfully", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusOK)

		body := response.Body.String()

		parsedUrl, err := url.Parse(body)
		if err != nil {
			t.Fatalf("Invalid download url %q returned", body)
		}

		query, err := url.ParseQuery(parsedUrl.RawQuery)
		if err != nil {
			t.Fatalf("Query string cannot be parsed %q", parsedUrl.RawQuery)
		}

		token := query.Get("token")
		if token == "" {
			t.Fatalf("Token string is missing, body %q", body)
		}

		assertFileUploadedProperly(t, mgr, token, fileContent)
	})
	t.Run("throws error for invalid request method", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodGet, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusMethodNotAllowed)
		assertResponseBody(t, response, MSG_ERR_INVALID_REQUEST_METHOD+"\n")
	})
	t.Run("throws error for invalid file field name", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "invalid", fileName, fileContent)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, MSG_ERR_CANNOT_READ_FILE+"\n")
	})
}

func TestDownload(t *testing.T) {
	mgr := NewStubFileManager()
	server := NewFileServer(mgr)

	buildDownloadUrl := func(token string) string {
		return fmt.Sprintf("%v?token=%v", DOWNLOAD_URL, url.QueryEscape(token))
	}

	t.Run("downloads successfully", func(t *testing.T) {
		fileName := "manual.txt"
		fileContent := "test content"

		buf := &bytes.Buffer{}
		buf.WriteString(fileContent)

		token, _ := mgr.SaveFile(fileName, buf)

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(token), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, fileContent)
		assertResponseFileHeaders(t, response, fileName, fileContent)
	})
	t.Run("returns error if filename query param is missing", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, DOWNLOAD_URL, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, MSG_ERR_MISSING_QUERY_PARAM+"\n")
	})
	t.Run("returns 404 if file not found", func(t *testing.T) {
		fileName := "nonexistent.txt"

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(fileName), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseStatus(t, response, http.StatusNotFound)
	})
}

// ------
// helper funcs
// ------

func createFileUploadRequest(method, fieldName, fileName, content string) *http.Request {
	buffer := bytes.Buffer{}
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	part, _ := writer.CreateFormFile(fieldName, fileName)
	fmt.Fprint(part, content)

	request := httptest.NewRequest(method, UPLOAD_URL, &buffer)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request
}

// ------
// asserts
// ------

func assertResponseStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()

	if response.Result().StatusCode != want {
		t.Errorf("Got status %d, but want %d", response.Result().StatusCode, want)
	}
}

func assertResponseBody(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()

	if response.Body.String() != want {
		t.Errorf("Got response %q, but want %q", response.Body.String(), want)
	}
}

func assertResponseHeader(t testing.TB, response *httptest.ResponseRecorder, header string, want []string) {
	t.Helper()

	values := response.Header().Values(header)

	if !reflect.DeepEqual(values, want) {
		t.Errorf("Got Header values %q, but want %q", values, want)
	}
}

func assertResponseFileHeaders(t testing.TB, response *httptest.ResponseRecorder, fileName, fileContent string) {
	t.Helper()

	contentLength := fmt.Sprint(len(fileContent))

	assertResponseHeader(t, response, "Content-Type", []string{"text/plain; charset=utf-8"})
	assertResponseHeader(t, response, "Content-Disposition", []string{"attachment; filename=" + fileName})
	assertResponseHeader(t, response, "Content-Length", []string{contentLength})
}

func assertFileUploadedProperly(t testing.TB, mgr *StubFileManager, token string, fileContent string) {
	t.Helper()

	loaded, ok := mgr.data[token]
	if !ok {
		t.Fatalf("Want %q to be in data, but it's not", token)
	}

	if loaded.content != fileContent {
		t.Fatalf("Got content %q, want %q", loaded.content, fileContent)
	}
}
