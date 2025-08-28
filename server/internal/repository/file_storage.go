package repository

import (
	"os"
	"path/filepath"

	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
)

// Ограничиваю доступ к файлам пользовательских данных
const perm os.FileMode = 0600

// FileStorage for manage file data
type FileStorage struct {
	Path     string
	FileName string
	file     *os.File
}

// FileStorageRepo for manage repository of files
type FileStorageRepo struct {
	fileStoragePath string
}

// InitFileStorage - init path to filestorage
func GetFileStorageRepo(storagePath string) (*FileStorageRepo, error) {
	var fileStorageR FileStorageRepo
	fileStorageR.fileStoragePath = storagePath
	// err := os.MkdirAll(storagePath, 0666)
	// if err != nil {
	// 	return nil, err
	// }
	return &fileStorageR, nil
}

// GetFileStorage return prepared filestorage fo filename
func (s *FileStorageRepo) GetFileStorage(fileName string) (FileStoragerInterface, error) {
	return &FileStorage{Path: s.fileStoragePath, FileName: fileName}, nil
}

func (r *FileStorage) getFilePath() (string, error) {
	return filepath.Abs(filepath.Join(r.Path, r.FileName))
}

// OpenWrite open file for write
func (r *FileStorage) OpenWrite() (*os.File, error) {
	var err error
	if r.file != nil {
		err = r.file.Close()
		r.file = nil
	}
	if err != nil {
		return nil, err
	}
	fpath, err := r.getFilePath()
	if err != nil {
		return nil, err
	}
	logger.Log.Debug("OpenWrite", zap.String("file path", fpath))
	file, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, perm)
	if err == nil {
		r.file = file
	} else {
		logger.Log.Error("OpenWrite", zap.Error(err))
	}
	return file, err
}

// OpenRead open file for read
func (r *FileStorage) OpenRead() (*os.File, error) {
	var err error
	if r.file != nil {
		err = r.file.Close()
		r.file = nil
	}
	if err != nil {
		return nil, err
	}
	fpath, err := r.getFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(fpath)
	if err == nil {
		r.file = file
	}
	return file, err
}

// Write write chunk  to file
func (r *FileStorage) Write(chunck []byte) (int, error) {
	return r.file.Write(chunck)
}

// Read data from file to  b
func (r *FileStorage) Read(b []byte) (int, error) {
	return r.file.Read(b)
}

// Close file
func (r *FileStorage) Close() error {
	if r.file == nil {
		return nil
	}
	return r.file.Close()
}
