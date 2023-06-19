package main

import (
	"io"
	"mime"
	"os"
	"path/filepath"
)

type UploadedFile struct {
	file      io.ReadCloser
	name      string
	size      int64
	mime_type string
}

type Storage interface {
	saveFile(fileName string, content io.Reader) error
	loadFile(fileName string) (uploaded UploadedFile, err error)
}

type FileSystemStorage struct {
	uploadDir string
}

func (f *FileSystemStorage) saveFile(fileName string, source io.Reader) error {
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

func (f *FileSystemStorage) loadFile(fileName string) (upload UploadedFile, err error) {
	path := f.buildPath(fileName)

	file, err := os.Open(path)

	if err != nil {
		return UploadedFile{}, err
	}

	stat, err := file.Stat()
	if err != nil {
		return UploadedFile{}, err
	}

	mime_type := mime.TypeByExtension(filepath.Ext(path))

	upload = UploadedFile{
		file:      file,
		name:      fileName,
		size:      stat.Size(),
		mime_type: mime_type,
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
