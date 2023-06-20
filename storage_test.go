package main

import (
	"bytes"
	"errors"
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
		NewFileSystemStoage(TMP_DIR)
		assertPathExists(t, TMP_DIR)
	})
	t.Run("does nothing if path already exists", func(t *testing.T) {
		defer setupTest()()

		os.MkdirAll(TMP_DIR, 0755)
		NewFileSystemStoage(TMP_DIR)
		assertPathExists(t, TMP_DIR)
	})
	t.Run("saves a file to the upload directory", func(t *testing.T) {
		defer setupTest()()

		fileName := "example.txt"
		fileContent := "content"

		storage := NewFileSystemStoage(TMP_DIR)

		buff := &bytes.Buffer{}
		buff.WriteString(fileContent)

		storage.saveFile(fileName, buff)

		assertFileSaved(t, fileName, fileContent)

		os.Remove(filepath.Join(TMP_DIR, fileName))
	})
	t.Run("loads file from the upload directory", func(t *testing.T) {
		defer setupTest()()

		fileName := "example.txt"
		fileContent := "test content"

		storage := NewFileSystemStoage(TMP_DIR)

		buff := &bytes.Buffer{}
		buff.WriteString(fileContent)

		storage.saveFile(fileName, buff)

		uploadedFile, _ := storage.loadFile(fileName)

		if uploadedFile.name != fileName {
			t.Errorf("Got name %q, want %q", uploadedFile.name, fileName)
		}

		fileSize := int64(len(fileContent))

		if uploadedFile.size != fileSize {
			t.Errorf("Got size %d, want %d", uploadedFile.size, fileSize)
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
