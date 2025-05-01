package file

import (
	"bin/interfaces"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type FileImpl struct{}

func NewFile() interfaces.File {
	return &FileImpl{}
}

func (*FileImpl) Read(source string) (data []byte, contentType string, err error) {

	data, err = os.ReadFile(source)
	if err != nil {
		return nil, "", fmt.Errorf("os.ReadFile error: %w", err)
	}

	contentType = http.DetectContentType(data)

	if ext := filepath.Ext(source); ext != "" {
		if mt := mime.TypeByExtension(ext); mt != "" {
			contentType = mt
		}
	}

	return data, contentType, nil
}

func (*FileImpl) Write(content []byte, name string) error {
	file, err := os.Create(name)
	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v", err)
		}
	}(file)

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}

var _ interfaces.File = (*FileImpl)(nil)
