package storage

import (
	"bin/interfaces"
	"encoding/json"
	"fmt"
	"strings"
)

type StorageImpl struct {
	fileClient interfaces.File
}
type BinIndex struct {
	Map map[string]string `json:"map"` // name → id
}

func NewStorage(f interfaces.File) interfaces.Storage {
	s := &StorageImpl{fileClient: f}

	_, _, err := s.fileClient.Read("b-bin-index.json")
	if err != nil {
		index := BinIndex{Map: map[string]string{}}
		data, _ := json.MarshalIndent(index, "", "  ")
		_ = s.fileClient.Write(data, "b-bin-index.json")
	}

	return s
}
func (s *StorageImpl) AddToIndex(name string, id string) error {
	data, _, err := s.fileClient.Read("b-bin-index.json")
	var index BinIndex

	if err == nil {
		_ = json.Unmarshal(data, &index)
	} else {
		index = BinIndex{Map: map[string]string{}}
	}

	index.Map[name] = id

	out, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return s.fileClient.Write(out, "b-bin-index.json")
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

func (s *StorageImpl) RemoveFromIndex(id string) error {
	data, _, err := s.fileClient.Read("b-bin-index.json")
	if err != nil {
		return err
	}

	var index BinIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return err
	}

	for name, savedId := range index.Map {
		if savedId == id {
			delete(index.Map, name)
			break
		}
	}

	out, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return s.fileClient.Write(out, "b-bin-index.json")
}

func (s *StorageImpl) GetIndex() (map[string]string, error) {
	data, _, err := s.fileClient.Read("b-bin-index.json")
	if err != nil {
		return nil, err
	}

	var index BinIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	return index.Map, nil
}
