package repository

import (
	"os"
	"path/filepath"

	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
)

const suffix = "fstorage"

type FileStorage struct {
	Path     string
	FileName string
	file     *os.File
}

type FileStorageRepo struct {
	fileStoragePath string
}

// InitFileStorage - init path to filestorage
func GetFileStorageRepo(storagePath string) (*FileStorageRepo, error) {
	var fileStorageR FileStorageRepo
	// if storagePath == "" {
	// 	dir, err := os.Getwd()
	// 	if err != nil {
	// 		return &fileStorageR, err
	// 	}
	// 	storagePath = filepath.Join(dir, suffix)
	// }
	// _, err := os.Stat(storagePath)
	// if errors.Is(err, fs.ErrNotExist) {
	// 	err = os.MkdirAll(storagePath, 0666)
	// }
	// if err != nil {
	// 	return &fileStorageR, err
	// }

	fileStorageR.fileStoragePath = storagePath
	return &fileStorageR, nil
}

func (s *FileStorageRepo) GetFileStorage(fileName string) (FileStoragerInterface, error) {
	return &FileStorage{Path: s.fileStoragePath, FileName: fileName}, nil
}

func (r *FileStorage) getFilePath() string {
	// return r.FileName
	return filepath.Join(r.Path, r.FileName)
}

func (r *FileStorage) OpenWrite() (*os.File, error) {
	var err error
	if r.file != nil {
		err = r.file.Close()
		r.file = nil
	}
	if err != nil {
		return nil, err
	}
	fpath := r.getFilePath()
	logger.Log.Debug("OpenWrite", zap.String("file path", fpath))
	file, err := os.OpenFile(fpath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
	if err == nil {
		r.file = file
	} else {
		logger.Log.Error("OpenWrite", zap.Error(err))
	}
	return file, err
}

func (r *FileStorage) OpenRead() (*os.File, error) {
	var err error
	if r.file != nil {
		err = r.file.Close()
		r.file = nil
	}
	if err != nil {
		return nil, err
	}
	fpath := r.getFilePath()
	file, err := os.Open(fpath)
	if err == nil {
		r.file = file
	}
	return file, err
}

func (r *FileStorage) Write(chunck []byte) (int, error) {
	return r.file.Write(chunck)
}

func (r *FileStorage) Read(b []byte) (int, error) {
	return r.file.Read(b)
}

func (r *FileStorage) Close() error {
	if r.file == nil {
		return nil
	}
	return r.file.Close()
}
