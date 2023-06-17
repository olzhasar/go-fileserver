package main

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	UPLOAD_DIR = "tmp"
}

func setupTest() func() {
	checkUploadDir()

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
		assertResponseBody(t, response, "File uploaded successfully")

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
		assertResponseBody(t, response, MSG_INVALID_REQUEST_METHOD+"\n")

		assertFileDoesNotExist(t, fileName)
	})
	t.Run("throws error for invalid file field name", func(t *testing.T) {
		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "invalid", fileName, fileContent)
		response := httptest.NewRecorder()

		uploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, MSG_CANNOT_READ_FILE+"\n")

		assertFileDoesNotExist(t, fileName)
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
