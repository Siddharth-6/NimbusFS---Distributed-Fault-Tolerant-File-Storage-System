package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BasePath string
}

func NewLocalStorage(path string) *LocalStorage {

	os.MkdirAll(path, 0755)

	return &LocalStorage{
		BasePath: path,
	}
}

func (l *LocalStorage) SaveChunk(id string, data []byte) error {

	path := filepath.Join(l.BasePath, id)

	return os.WriteFile(path, data, 0644)
}

func (l *LocalStorage) OpenChunk(id string) (io.ReadCloser, error) {

	path := filepath.Join(l.BasePath, id)

	return os.Open(path)
}

func (l *LocalStorage) DeleteChunk(id string) error {

	path := filepath.Join(l.BasePath, id)

	return os.Remove(path)
}
