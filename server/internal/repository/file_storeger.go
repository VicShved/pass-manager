package repository

import (
	// "io"
	"os"
)

type FileStoragerRepoInterface interface {
	GetFileStorage(fileName string) (FileStoragerInterface, error)
}

// FileStoragerInterface interface
type FileStoragerInterface interface {
	// Open(fileName string, flag int, perm os.FileMode) (*os.File, error)
	OpenWrite() (*os.File, error)
	OpenRead() (*os.File, error)
	Write(chunck []byte) (int, error)
	Read(b []byte) (int, error)
	// ToFile(r io.Reader) (int, error)
	Close() error
}
