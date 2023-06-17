package main

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func init() {
	UPLOAD_DIR = "tmp"
	storage = &FileSystemStorage{}
}

func setupTest() func() {
	createDirIfNotExists(UPLOAD_DIR)

	return func() {
		os.RemoveAll(UPLOAD_DIR)
	}
}

func TestUpload(t *testing.T) {
	defer setupTest()()

	t.Run("uploads successfully", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		uploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, MSG_UPLOAD_SUCCESS)

		assertFileSaved(t, fileName, fileContent)
		deleteUploadedFile(t, fileName)
		assertFileDoesNotExist(t, fileName)
	})
	t.Run("throws error for invalid request method", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodGet, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		uploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusMethodNotAllowed)
		assertResponseBody(t, response, MSG_ERR_INVALID_REQUEST_METHOD+"\n")

		assertFileDoesNotExist(t, fileName)
	})
	t.Run("throws error for invalid file field name", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "invalid", fileName, fileContent)
		response := httptest.NewRecorder()

		uploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, MSG_ERR_CANNOT_READ_FILE+"\n")

		assertFileDoesNotExist(t, fileName)
	})
}

func TestDownload(t *testing.T) {
	defer setupTest()()

	buildDownloadUrl := func(fileName string) string {
		return fmt.Sprintf("%v?filename=%v", DOWNLOAD_URL, url.QueryEscape(fileName))
	}

	t.Run("downloads successfully", func(t *testing.T) {
		fileName := "manual.txt"
		fileContent := "test content"

		createUploadedFile(fileName, fileContent)
		assertFileSaved(t, fileName, fileContent)

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(fileName), &bytes.Buffer{})
		response := httptest.NewRecorder()

		downloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, fileContent)
		assertResponseFileHeaders(t, response, fileName, fileContent)
	})
	t.Run("returns error if filename query param is missing", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, DOWNLOAD_URL, &bytes.Buffer{})
		response := httptest.NewRecorder()

		downloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, MSG_ERR_MISSING_QUERY_PARAM+"\n")
	})
	t.Run("returns 404 if file not found", func(t *testing.T) {
		fileName := "nonexistent.txt"

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(fileName), &bytes.Buffer{})
		response := httptest.NewRecorder()

		downloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusNotFound)
	})
}

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

func createUploadedFile(fileName string, fileContent string) {
	path := filepath.Join(UPLOAD_DIR, fileName)
	os.WriteFile(path, []byte(fileContent), 0644)
}

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

func assertFileSaved(t testing.TB, fileName, want string) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(UPLOAD_DIR, fileName))
	if err != nil {
		t.Fatalf("Error opening file:\n%q", err)
	}

	content := string(data)

	if content != want {
		t.Errorf("Got %q, but want %q", data, want)
	}
}

func assertFileDoesNotExist(t testing.TB, fileName string) {
	_, err := os.Stat(filepath.Join(UPLOAD_DIR, fileName))

	if err == nil {
		t.Fatal("expected file to not exist, but it does")
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Unexpected error:\n%q", err)
	}
}

func deleteUploadedFile(t testing.TB, fileName string) {
	err := os.Remove(filepath.Join(UPLOAD_DIR, fileName))

	if err != nil {
		t.Fatalf("Error deleting file:\n%q", err)
	}
}
