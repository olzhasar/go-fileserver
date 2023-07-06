package manager

import (
	"bytes"
	"testing"

	"github.com/olzhasar/go-fileserver/registry"
	"github.com/olzhasar/go-fileserver/storages"
)

func TestSaveFile(t *testing.T) {
	reg := registry.NewInMemoryRegistry()
	storage := storages.NewInMemoryStorage()
	mgr := &FileManager{registry: reg, storage: storage}

	fileName := "example.txt"
	fileContent := "test content"

	buf := &bytes.Buffer{}
	buf.WriteString(fileContent)

	token, err := mgr.SaveFile(fileName, buf)

	if err != nil {
		t.Fatalf("Expected no error, got %q", err)
	}

	savedName, ok := reg.Get(token)

	if !ok {
		t.Fatalf("Want token %q to be in registry, but it's not", token)
	}

	if savedName != fileName {
		t.Fatalf("Got filename %q, want %q", savedName, fileName)
	}

	upload, err := storage.LoadFile(fileName)

	if err != nil {
		t.Fatalf("Error loading file from storage:\n%v", err)
	}

	if upload.Name != fileName {
		t.Fatalf("Got upload name %q, want %q", upload.Name, fileName)
	}
}

func TestLoadFile(t *testing.T) {
	reg := registry.NewInMemoryRegistry()
	storage := storages.NewInMemoryStorage()
	mgr := &FileManager{registry: reg, storage: storage}

	t.Run("loads file from storage", func(t *testing.T) {
		fileName := "example.txt"
		fileContent := "test content"

		buf := &bytes.Buffer{}
		buf.WriteString(fileContent)

		token, _ := mgr.SaveFile(fileName, buf)

		upload, err := mgr.LoadFile(token)

		if err != nil {
			t.Fatalf("Expected no error, got %q", err)
		}

		if upload.Name != fileName {
			t.Fatalf("Got filename %q, want %q", upload.Name, fileName)
		}
	})

	t.Run("throws error for unexisting file", func(t *testing.T) {
		_, err := mgr.LoadFile("123456")

		if err == nil {
			t.Fatal("Got nil, want error")
		}
	})
}
