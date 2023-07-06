package manager

import (
	"errors"
	"io"

	"github.com/olzhasar/go-fileserver/registry"
	"github.com/olzhasar/go-fileserver/storages"
)

type SaverLoader interface {
	SaveFile(fileName string, content io.Reader) (token string, err error)
	LoadFile(token string) (upload storages.UploadedFile, err error)
}

type FileManager struct {
	registry registry.Registry
	storage  storages.Storage
}

func NewFileManager(r registry.Registry, s storages.Storage) *FileManager {
	return &FileManager{r, s}
}

func (f *FileManager) SaveFile(fileName string, content io.Reader) (string, error) {
	token, err := registry.RecordFile(f.registry, fileName, registry.GenerateUniqueToken)
	if err != nil {
		return "", err
	}

	err = f.storage.SaveFile(fileName, content)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (f *FileManager) LoadFile(token string) (storages.UploadedFile, error) {
	fileName, ok := f.registry.Get(token)
	if !ok {
		return storages.UploadedFile{}, errors.New("Invalid token")
	}

	upload, err := f.storage.LoadFile(fileName)

	if err != nil {
		return storages.UploadedFile{}, err
	}

	return upload, nil
}
