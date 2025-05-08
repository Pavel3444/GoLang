package storage

import (
	"bin/interfaces"
	"fmt"
	"strings"
)

type StorageImpl struct {
	fileClient interfaces.File
}

func NewStorage(f interfaces.File) interfaces.Storage {
	return &StorageImpl{fileClient: f}
}

func (s *StorageImpl) Save(data []byte, name string) error {
	return s.fileClient.Write(data, name)
}

func (s *StorageImpl) Read(name string) ([]byte, error) {
	data, fileType, err := s.fileClient.Read(name)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(fileType, "application/json") {
		return nil, fmt.Errorf("expected application/json, got %q", fileType)
	}
	return data, nil
}

var _ interfaces.Storage = (*StorageImpl)(nil)
