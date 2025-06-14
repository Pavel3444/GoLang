package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type JSONBinClient struct {
	ApiKey  string
	BaseURL string
}

type CreatedBin struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewJSONBinClient() *JSONBinClient {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, relying on existing environment variables")
	}
	apiKey := os.Getenv("KEY")
	if apiKey == "" {
		log.Fatal("API key is missing. Please set KEY in the environment or .env file.")
	}
	return &JSONBinClient{
		ApiKey:  apiKey,
		BaseURL: "https://api.jsonbin.io/v3/b",
	}
}

func (c *JSONBinClient) CreateBin(filePath, name string) (*CreatedBin, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %q: %w", filePath, err)
	}

	var record interface{}
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("не удалось распарсить JSON в файле %q: %w", filePath, err)
	}

	payload := struct {
		Name   string      `json:"name"`
		Record interface{} `json:"record"`
	}{
		Name:   name,
		Record: record,
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("не удалось собрать JSON для запроса: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.BaseURL, buf)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать HTTP-запрос: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе к серверу: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf(
				"неожиданный HTTP-статус %d и ошибка при чтении тела ответа: %w",
				resp.StatusCode, err,
			)
		}
		return nil, fmt.Errorf(
			"неожиданный HTTP-статус %d: %s",
			resp.StatusCode, bytes.TrimSpace(body),
		)
	}

	var result struct {
		Metadata struct {
			ID string `json:"id"`
		} `json:"metadata"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("не удалось распарсить ответ сервера: %w", err)
	}

	return &CreatedBin{
		Id:   result.Metadata.ID,
		Name: name,
	}, nil
}

func (c *JSONBinClient) GetBin(id string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/latest", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать HTTP-запрос: %w", err)
	}
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении HTTP-запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf(
				"неожиданный HTTP-статус %d и ошибка при чтении тела ответа: %w",
				resp.StatusCode, readErr,
			)
		}
		return nil, fmt.Errorf(
			"неожиданный HTTP-статус %d: %s",
			resp.StatusCode, strings.TrimSpace(string(body)),
		)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать тело ответа: %w", err)
	}
	return data, nil
}
func (c *JSONBinClient) UpdateBin(filePath, id string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл %q: %w", filePath, err)
	}

	var tmp interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("невалидный JSON в файле %q: %w", filePath, err)
	}

	url := fmt.Sprintf("%s/%s", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("не удалось создать запрос: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf(
				"неожиданный HTTP-статус %d и ошибка при чтении тела: %w",
				resp.StatusCode, readErr,
			)
		}
		return fmt.Errorf(
			"не удалось обновить бин (status %d): %s",
			resp.StatusCode, strings.TrimSpace(string(body)),
		)
	}

	return nil
}

func (c *JSONBinClient) DeleteBin(id string) error {
	url := fmt.Sprintf("%s/%s", c.BaseURL, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("не удалось создать HTTP-запрос для удаления: %w", err)
	}
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении HTTP-запроса на удаление: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf(
				"неожиданный HTTP-статус %d и ошибка при чтении тела ответа: %w",
				resp.StatusCode, readErr,
			)
		}
		return fmt.Errorf(
			"не удалось удалить бин (status %d): %s",
			resp.StatusCode, strings.TrimSpace(string(body)),
		)
	}

	return nil
}
