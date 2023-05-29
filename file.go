package main

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

const UPLOAD_DIR = "uploads"

func getUploadFilePath(filename string) string {
	return filepath.Join(UPLOAD_DIR, filename)
}

func saveFile(filename string, f *multipart.File) error {
	newFilePath := getUploadFilePath(filename)
	newFile, err := os.Create(newFilePath)

	if err != nil {
		return err
	}

	defer newFile.Close()

	_, err = io.Copy(newFile, *f)
	if err != nil {
		return err
	}

	return nil
}
