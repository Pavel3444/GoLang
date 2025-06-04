package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	_ = godotenv.Load(".env")
	return &JSONBinClient{
		ApiKey:  os.Getenv("KEY"),
		BaseURL: "https://api.jsonbin.io/v3/b",
	}
}

func (c *JSONBinClient) CreateBin(filePath string, name string) (*CreatedBin, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	var content interface{}
	if err := json.Unmarshal(file, &content); err != nil {
		return nil, fmt.Errorf("не удалось распарсить JSON: %w", err)
	}

	binData := map[string]interface{}{
		"name":   name,
		"record": content,
	}

	jsonData, _ := json.Marshal(binData)

	req, _ := http.NewRequest("POST", c.BaseURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Metadata struct {
			ID string `json:"id"`
		} `json:"metadata"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &CreatedBin{
		Id:   result.Metadata.ID,
		Name: name,
	}, nil
}

func (c *JSONBinClient) GetBin(id string) ([]byte, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s/latest", c.BaseURL, id), nil)
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *JSONBinClient) UpdateBin(filePath string, id string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл: %w", err)
	}

	var content interface{}
	if err := json.Unmarshal(file, &content); err != nil {
		return fmt.Errorf("невалидный JSON: %w", err)
	}

	jsonData, err := json.Marshal(content)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/%s", c.BaseURL, id), bytes.NewReader(jsonData))
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при обновлении бин: %s", string(body))
	}

	return nil
}

func (c *JSONBinClient) DeleteBin(id string) error {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", c.BaseURL, id), nil)
	req.Header.Set("X-Master-Key", c.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("ошибка при удалении бин: %s", string(body))
	}
	return nil
}
