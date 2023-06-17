package main

import (
	"io"
	"os"
	"path/filepath"
)

func getUploadFilePath(fileName string) string {
	return filepath.Join(UPLOAD_DIR, fileName)
}

func saveFile(fileName string, content io.Reader) error {
	newFilePath := getUploadFilePath(fileName)
	newFile, err := os.Create(newFilePath)

	if err != nil {
		return err
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, content)
	if err != nil {
		return err
	}

	return nil
}

func createDirIfNotExists(path string) error {
	return os.MkdirAll(path, 0755)
}
