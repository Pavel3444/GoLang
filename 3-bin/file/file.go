package file

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func ReadFile(source string) (data []byte, contentType string, err error) {
	u, parseErr := url.Parse(source)
	if parseErr == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(source)
		if err != nil {
			return nil, "", fmt.Errorf("HTTP GET error: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, "", fmt.Errorf("HTTP GET returned status %s", resp.Status)
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, "", fmt.Errorf("reading HTTP body: %w", err)
		}
	} else {
		data, err = os.ReadFile(source)
		if err != nil {
			return nil, "", fmt.Errorf("os.ReadFile error: %w", err)
		}
	}

	contentType = http.DetectContentType(data)
	if u == nil || u.Scheme == "" {
		if ext := filepath.Ext(source); ext != "" {
			if mt := mime.TypeByExtension(ext); mt != "" {
				contentType = mt
			}
		}
	}

	return data, contentType, nil
}

func WriteFile(content []byte, name string) error {
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
