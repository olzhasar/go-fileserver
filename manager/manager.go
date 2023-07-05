package manager

import (
	"github.com/olzhasar/go-fileserver/storages"
	"io"
)

type FileManager interface {
	SaveFile(fileName string, content io.Reader) (token string, err error)
	LoadFile(token string) (upload storages.UploadedFile, err error)
}
