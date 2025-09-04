package storage

import (
	"encoding/json"
	"handyhub-email-svc/internal/config"
	"handyhub-email-svc/internal/models"
	"os"
	"path/filepath"
)

type FileStorage struct {
	config config.FileStorageConfig
	file   *os.File
}

func NewFileStorage(cfg config.FileStorageConfig) (*FileStorage, error) {

	if err := os.MkdirAll(filepath.Dir(cfg.Path), os.ModePerm); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(cfg.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.Info("File storage initialized at ", cfg.Path)

	return &FileStorage{
		config: cfg,
		file:   file,
	}, nil
}

func (fs *FileStorage) Store(emailLog *models.EmailLog) error {
	data, err := json.Marshal(emailLog)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = fs.file.Write(data)
	if err != nil {
		return err
	}
	log.Info("Email log entry stored in file")
	return nil
}

func (fs *FileStorage) Close() error {
	if err := fs.file.Close(); err != nil {
		return err
	}
	log.Info("File storage closed")
	return nil
}
