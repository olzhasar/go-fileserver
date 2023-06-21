package main_test

import (
	"bytes"
	"errors"
	"github.com/olzhasar/go-fileserver"
	"os"
	"path/filepath"
	"testing"
)

const TMP_DIR = "tmp"

func TestFileSystemStorage(t *testing.T) {
	setupTest := func() func() {
		return func() {
			os.RemoveAll(TMP_DIR)
		}
	}

	t.Run("creates a directory if it does not exist", func(t *testing.T) {
		defer setupTest()()

		assertPathDoesNotExist(t, TMP_DIR)
		main.NewFileSystemStoage(TMP_DIR)
		assertPathExists(t, TMP_DIR)
	})
	t.Run("does nothing if path already exists", func(t *testing.T) {
		defer setupTest()()

		os.MkdirAll(TMP_DIR, 0755)
		main.NewFileSystemStoage(TMP_DIR)
		assertPathExists(t, TMP_DIR)
	})
	t.Run("saves a file to the upload directory", func(t *testing.T) {
		defer setupTest()()

		fileName := "example.txt"
		fileContent := "content"
		buff := createContentBuffer(fileContent)

		storage := main.NewFileSystemStoage(TMP_DIR)

		storage.SaveFile(fileName, buff)

		assertFileSaved(t, fileName, fileContent)
	})
	t.Run("loads file from the upload directory", func(t *testing.T) {
		defer setupTest()()

		fileName := "example.txt"
		fileContent := "test content"

		storage := main.NewFileSystemStoage(TMP_DIR)

		buff := &bytes.Buffer{}
		buff.WriteString(fileContent)

		storage.SaveFile(fileName, buff)

		uploadedFile, _ := storage.LoadFile(fileName)

		checkUploadedFile(t, uploadedFile, fileName, fileContent)
	})
}

func TestInMemoryStorage(t *testing.T) {
	t.Run("Saves file to memory", func(t *testing.T) {
		fileName := "example.txt"
		fileContent := "content"

		buff := createContentBuffer(fileContent)

		storage := main.NewInMemoryStorage()
		storage.SaveFile(fileName, buff)

		if val, ok := storage.Files[fileName]; !ok {
			t.Errorf("Expected %q to be in storage, but it wasn't", fileName)
		} else if val != fileContent {
			t.Errorf("got file content %q, want %q", val, fileContent)
		}
	})
	t.Run("Loads file from memory", func(t *testing.T) {
		fileName := "example.txt"
		fileContent := "content"

		buff := createContentBuffer(fileContent)

		storage := main.NewInMemoryStorage()
		storage.SaveFile(fileName, buff)

		uploaded, err := storage.LoadFile(fileName)

		if err != nil {
			t.Fatalf("Expected no error, but got %q", err)
		}

		checkUploadedFile(t, uploaded, fileName, fileContent)

		defer uploaded.File.Close()
	})
	t.Run("Throws error when loading missing file", func(t *testing.T) {
		storage := main.NewInMemoryStorage()
		_, err := storage.LoadFile("nonexisting.txt")

		if err == nil {
			t.Fatal("Expected error, but did not get one")
		}
	})
	t.Run("Clear deletes everything from map", func(t *testing.T) {
		storage := main.NewInMemoryStorage()

		buff := createContentBuffer("test")
		storage.SaveFile("test.txt", buff)

		storage.Clear()

		if len(storage.Files) != 0 {
			t.Fatalf("Expected storage files map to be clear, got %d entries", len(storage.Files))
		}
	})
}

func createContentBuffer(content string) *bytes.Buffer {
	buff := &bytes.Buffer{}
	buff.WriteString(content)
	return buff
}

func checkUploadedFile(t testing.TB, uploadedFile main.UploadedFile, name, content string) {
	t.Helper()

	if uploadedFile.Name != name {
		t.Errorf("Got name %q, want %q", uploadedFile.Name, name)
	}

	fileSize := int64(len(content))

	if uploadedFile.Size != fileSize {
		t.Errorf("Got size %d, want %d", uploadedFile.Size, fileSize)
	}
}

func assertPathDoesNotExist(t testing.TB, path string) {
	t.Helper()

	_, err := os.Stat(path)

	if err == nil {
		t.Fatalf("expected path %q to not exist, but it does", path)
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Unexpected error:\n%q", err)
	}
}

func assertPathExists(t testing.TB, path string) {
	t.Helper()

	_, err := os.Stat(filepath.Join(path))

	if err != nil {
		t.Fatalf("expected path %q to exist, but it does not\n%q", path, err)
	}
}

func assertFileSaved(t testing.TB, fileName, want string) {
	t.Helper()

	data, err := os.ReadFile(filepath.Join(TMP_DIR, fileName))
	if err != nil {
		t.Fatalf("Error opening file:\n%q", err)
	}

	content := string(data)

	if content != want {
		t.Errorf("Got %q, but want %q", data, want)
	}
}
