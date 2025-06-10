package api

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found, relying on existing environment variables")
	}
}

func createTempJSONFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.json")
	if err != nil {
		t.Fatalf("не удалось создать temp-файл: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })

	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("не удалось записать в temp-файл: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("не удалось закрыть temp-файл: %v", err)
	}
	return f.Name()
}

func createTestBin(t *testing.T, client *JSONBinClient, content, name string) *CreatedBin {
	t.Helper()
	path := createTempJSONFile(t, content)
	bin, err := client.CreateBin(path, name)
	require.NoError(t, err, "CreateBin(%s) failed", name)

	t.Cleanup(func() {
		if err := client.DeleteBin(bin.Id); err != nil {
			t.Errorf("DeleteBin(%s) failed: %v", bin.Id, err)
		}
	})
	return bin
}

func TestCreateBin(t *testing.T) {
	client := NewJSONBinClient()
	bin := createTestBin(t, client, `{"example":"data"}`, "test-create")
	require.NotEmpty(t, bin.Id, "ожидался непустой ID")
}

func TestGetBin(t *testing.T) {
	client := NewJSONBinClient()
	bin := createTestBin(t, client, `{"foo":"bar"}`, "test-get")

	data, err := client.GetBin(bin.Id)
	require.NoError(t, err, "GetBin(%s) should succeed", bin.Id)
	assert.Contains(t, string(data), `"foo":"bar"`)
}

func TestUpdateBin(t *testing.T) {
	client := NewJSONBinClient()
	bin := createTestBin(t, client, `{"a":1}`, "test-update")

	updatePath := createTempJSONFile(t, `{"a":2}`)
	require.NoError(t, client.UpdateBin(updatePath, bin.Id), "UpdateBin should not error")

	data, err := client.GetBin(bin.Id)
	require.NoError(t, err)
	assert.Contains(t, string(data), `"a":2`)
}

func TestDeleteBin(t *testing.T) {
	client := NewJSONBinClient()
	bin := createTestBin(t, client, `{"x":"y"}`, "test-delete")

	require.NoError(t, client.DeleteBin(bin.Id), "DeleteBin should not error")

	_, err := client.GetBin(bin.Id)
	assert.Error(t, err, "после удаления GetBin должен вернуть ошибку")
}
