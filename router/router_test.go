package router_test

import (
	"bytes"
	"fmt"
	"github.com/olzhasar/go-fileserver/router"
	"github.com/olzhasar/go-fileserver/storages"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func setupTest(s *storages.InMemoryStorage) func() {
	return func() {
		s.Clear()
	}
}

func TestUpload(t *testing.T) {
	storage := storages.NewInMemoryStorage()
	rt := router.NewRouter(storage)

	t.Run("uploads successfully", func(t *testing.T) {
		defer setupTest(storage)()

		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		rt.UploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, router.MSG_UPLOAD_SUCCESS)

		uploaded, err := storage.LoadFile(fileName)
		if err != nil {
			t.Fatalf("Got error %q while loading uploaded file from storage", err)
		}

		if uploaded.Name != fileName {
			t.Errorf("Got name %q, but want %q", uploaded.Name, fileName)
		}

		buff := &bytes.Buffer{}
		io.Copy(buff, uploaded.File)
		uploadedContent := buff.String()

		if uploadedContent != fileContent {
			t.Errorf("Got content %q, but want %q", uploadedContent, fileContent)
		}

		size := int64(len(fileContent))
		if uploaded.Size != size {
			t.Errorf("Got size %q, but want %q", uploaded.Size, size)
		}
	})
	t.Run("throws error for invalid request method", func(t *testing.T) {
		defer setupTest(storage)()

		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodGet, "file", fileName, fileContent)
		response := httptest.NewRecorder()

		rt.UploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusMethodNotAllowed)
		assertResponseBody(t, response, router.MSG_ERR_INVALID_REQUEST_METHOD+"\n")
		assertStorageIsEmpty(t, storage)
	})
	t.Run("throws error for invalid file field name", func(t *testing.T) {
		defer setupTest(storage)()

		fileName := "test_file.txt"
		fileContent := "test content"

		request := createFileUploadRequest(http.MethodPost, "invalid", fileName, fileContent)
		response := httptest.NewRecorder()

		rt.UploadHandler(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, router.MSG_ERR_CANNOT_READ_FILE+"\n")
		assertStorageIsEmpty(t, storage)
	})
}

func TestDownload(t *testing.T) {
	storage := storages.NewInMemoryStorage()
	rt := router.NewRouter(storage)

	buildDownloadUrl := func(fileName string) string {
		return fmt.Sprintf("%v?filename=%v", router.DOWNLOAD_URL, url.QueryEscape(fileName))
	}

	t.Run("downloads successfully", func(t *testing.T) {
		defer setupTest(storage)()

		fileName := "manual.txt"
		fileContent := "test content"
		createUploadedFile(storage, fileName, fileContent)

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(fileName), &bytes.Buffer{})
		response := httptest.NewRecorder()

		rt.DownloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusOK)
		assertResponseBody(t, response, fileContent)
		assertResponseFileHeaders(t, response, fileName, fileContent)
	})
	t.Run("returns error if filename query param is missing", func(t *testing.T) {
		defer setupTest(storage)()

		request := httptest.NewRequest(http.MethodGet, router.DOWNLOAD_URL, &bytes.Buffer{})
		response := httptest.NewRecorder()

		rt.DownloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusBadRequest)
		assertResponseBody(t, response, router.MSG_ERR_MISSING_QUERY_PARAM+"\n")
	})
	t.Run("returns 404 if file not found", func(t *testing.T) {
		defer setupTest(storage)()

		fileName := "nonexistent.txt"

		request := httptest.NewRequest(http.MethodGet, buildDownloadUrl(fileName), &bytes.Buffer{})
		response := httptest.NewRecorder()

		rt.DownloadHandler(response, request)

		assertResponseStatus(t, response, http.StatusNotFound)
	})
}

func createFileUploadRequest(method, fieldName, fileName, content string) *http.Request {
	buffer := bytes.Buffer{}
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	part, _ := writer.CreateFormFile(fieldName, fileName)
	fmt.Fprint(part, content)

	request := httptest.NewRequest(method, router.UPLOAD_URL, &buffer)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request
}

func createUploadedFile(storage storages.Storage, fileName string, fileContent string) {
	buffer := &bytes.Buffer{}
	buffer.WriteString(fileContent)

	storage.SaveFile(fileName, buffer)
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

func assertStorageIsEmpty(t testing.TB, storage *storages.InMemoryStorage) {
	t.Helper()

	if len(storage.Files) != 0 {
		t.Errorf("Expected storage to be empty, but got %d entries", len(storage.Files))
	}
}
