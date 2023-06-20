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

		storage := main.NewFileSystemStoage(TMP_DIR)

		buff := &bytes.Buffer{}
		buff.WriteString(fileContent)

		storage.SaveFile(fileName, buff)

		assertFileSaved(t, fileName, fileContent)

		os.Remove(filepath.Join(TMP_DIR, fileName))
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

		if uploadedFile.Name != fileName {
			t.Errorf("Got name %q, want %q", uploadedFile.Name, fileName)
		}

		fileSize := int64(len(fileContent))

		if uploadedFile.Size != fileSize {
			t.Errorf("Got size %d, want %d", uploadedFile.Size, fileSize)
		}
	})
}

func assertPathDoesNotExist(t testing.TB, path string) {
	_, err := os.Stat(path)

	if err == nil {
		t.Fatalf("expected path %q to not exist, but it does", path)
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Unexpected error:\n%q", err)
	}
}

func assertPathExists(t testing.TB, path string) {
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
