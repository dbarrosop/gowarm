package core

import (
	"encoding/json"
	"os"
	"sync"
)

type FileStorage struct {
	f  *os.File
	mu *sync.Mutex
}

func NewFileStorage(f *os.File) *FileStorage {
	return &FileStorage{f, &sync.Mutex{}}
}

func (fs *FileStorage) Close() error {
	return fs.f.Close()
}

func (fs *FileStorage) Write(data interface{}) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	if err := fs.f.Truncate(0); err != nil {
		return err
	}
	if _, err := fs.f.Seek(0, 0); err != nil {
		return err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fs.f.Write(b)
	return err
}
