package storages

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
)

type UploadedFile struct {
	File io.ReadCloser
	Name string
	Size int64
}

func (u *UploadedFile) MimeTypeByExt() string {
	return mime.TypeByExtension(filepath.Ext(u.Name))
}

type Storage interface {
	SaveFile(fileName string, content io.Reader) error
	LoadFile(fileName string) (uploaded UploadedFile, err error)
}

type FileSystemStorage struct {
	uploadDir string
}

func (f *FileSystemStorage) SaveFile(fileName string, source io.Reader) error {
	newFilePath := f.buildPath(fileName)
	newFile, err := os.Create(newFilePath)

	if err != nil {
		return err
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, source)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileSystemStorage) LoadFile(fileName string) (upload UploadedFile, err error) {
	path := f.buildPath(fileName)

	file, err := os.Open(path)

	if err != nil {
		return UploadedFile{}, err
	}

	stat, err := file.Stat()
	if err != nil {
		return UploadedFile{}, err
	}

	upload = UploadedFile{
		File: file,
		Name: fileName,
		Size: stat.Size(),
	}

	return upload, nil
}

func (f *FileSystemStorage) buildPath(fileName string) string {
	return filepath.Join(f.uploadDir, fileName)
}

func NewFileSystemStoage(uploadDir string) Storage {
	os.MkdirAll(uploadDir, 0755)

	return &FileSystemStorage{uploadDir: uploadDir}
}

type InMemoryStorage struct {
	Files map[string]string
}

type InMemoryFile struct {
	Buffer io.Reader
}

func (i InMemoryFile) Read(p []byte) (int, error) {
	return i.Buffer.Read(p)
}

func (i InMemoryFile) Close() error {
	return nil
}

func (i *InMemoryStorage) SaveFile(fileName string, source io.Reader) error {
	buff := &bytes.Buffer{}
	io.Copy(buff, source)
	i.Files[fileName] = buff.String()
	return nil
}

func (i *InMemoryStorage) LoadFile(fileName string) (UploadedFile, error) {
	content, ok := i.Files[fileName]

	if !ok {
		return UploadedFile{}, errors.New(fmt.Sprintf("File %q not found in storage", fileName))
	}

	buff := &bytes.Buffer{}
	buff.WriteString(content)
	file := InMemoryFile{buff}

	size := int64(len(content))

	return UploadedFile{File: file, Name: fileName, Size: size}, nil
}

func (i *InMemoryStorage) Clear() {
	for k := range i.Files {
		delete(i.Files, k)
	}
}

func NewInMemoryStorage() *InMemoryStorage {
	files := make(map[string]string)
	return &InMemoryStorage{files}
}
