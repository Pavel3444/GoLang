package storage

import (
	"bin/bins"
	"bin/file"
	"encoding/json"
	"fmt"
	"strings"
)

func SaveBins(content []byte, name string) error {
	return file.WriteFile(content, name)
}

func ReadBins(name string, list *bins.BinList) ([]byte, error) {
	data, fileType, err := file.ReadFile(name)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(fileType, "application/json") {
		return nil, fmt.Errorf("expected 'application/json' but got '%s'", fileType)
	}
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("invalid JSON in %q: %w", name, err)
	}
	return data, nil
}
