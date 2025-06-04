package api

import (
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("не удалось загрузить .env: " + err.Error())
	}
}

func TestCreateBin(t *testing.T) {
	client := NewJSONBinClient()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`{"example": "data"}`)
	tmpFile.Close()

	bin, err := client.CreateBin(tmpFile.Name(), "test-bin")
	if err != nil {
		t.Fatalf("ошибка при создании бин: %v", err)
	}
	if bin.Id == "" {
		t.Error("ожидался непустой ID")
	}

	defer client.DeleteBin(bin.Id)
}

func TestGetBin(t *testing.T) {
	client := NewJSONBinClient()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`{"example": "data"}`)
	tmpFile.Close()

	bin, err := client.CreateBin(tmpFile.Name(), "get-test-bin")
	if err != nil {
		t.Fatalf("ошибка при создании бин: %v", err)
	}
	defer client.DeleteBin(bin.Id)

	_, err = client.GetBin(bin.Id)
	if err != nil {
		t.Errorf("ошибка при получении бин: %v", err)
	}
}

func TestUpdateBin(t *testing.T) {
	client := NewJSONBinClient()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.WriteString(`{"initial": "data"}`)
	tmpFile.Close()

	bin, err := client.CreateBin(tmpFile.Name(), "update-test-bin")
	if err != nil {
		t.Fatalf("ошибка при создании бин: %v", err)
	}
	defer client.DeleteBin(bin.Id)

	updateFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(updateFile.Name())
	updateFile.WriteString(`{"updated": "data"}`)
	updateFile.Close()

	err = client.UpdateBin(updateFile.Name(), bin.Id)
	if err != nil {
		t.Errorf("ошибка при обновлении бин: %v", err)
	}
}

func TestDeleteBin(t *testing.T) {
	client := NewJSONBinClient()

	tmpFile, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString(`{"to_be_deleted": "yes"}`)
	tmpFile.Close()

	bin, err := client.CreateBin(tmpFile.Name(), "delete-test-bin")
	if err != nil {
		t.Fatalf("ошибка при создании бин: %v", err)
	}

	err = client.DeleteBin(bin.Id)
	if err != nil {
		t.Errorf("ошибка при удалении бин: %v", err)
	}
}
