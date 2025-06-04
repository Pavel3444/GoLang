package main

import (
	"bin/api"
	"bytes"
	"encoding/json"

	//"bin/bins"
	"bin/file"
	"bin/storage"
	"flag"
	"fmt"
	"log"
)

func main() {
	fileClient := file.NewFile()
	store := storage.NewStorage(fileClient)
	//binRepo := bins.NewBinRepository(store, "bins.json")
	apiClient := api.NewJSONBinClient()

	createCmd := flag.Bool("create", false, "Create a new bin")
	updateCmd := flag.Bool("update", false, "Update an existing bin")
	deleteCmd := flag.Bool("delete", false, "Delete a bin")
	getCmd := flag.Bool("get", false, "Get a bin")
	listCmd := flag.Bool("list", false, "List all bins")

	fileFlag := flag.String("file", "", "Path to file with JSON")
	nameFlag := flag.String("name", "", "Name of the bin")
	idFlag := flag.String("id", "", "ID of the bin")

	flag.Parse()

	switch {
	case *createCmd:
		if *fileFlag == "" || *nameFlag == "" {
			log.Fatal("Для --create нужны --file и --name")
		}
		bin, err := apiClient.CreateBin(*fileFlag, *nameFlag)
		if err != nil {
			log.Fatalf("Ошибка при создании: %v", err)
		}
		fmt.Printf("Создан bin: %+v\n", bin)
		if err := store.AddToIndex(bin.Name, bin.Id); err != nil {
			log.Printf("Не удалось сохранить бин в индекс: %v\n", err)
		}
	case *updateCmd:
		if *fileFlag == "" || *idFlag == "" {
			log.Fatal("Для --create нужны --file и --id")
		}
		if err := apiClient.UpdateBin(*fileFlag, *idFlag); err != nil {
			log.Fatalf("Ошибка при обновлении: %v", err)
		}
		fmt.Println("Бин успешно обновлён.")
	case *getCmd:
		if *idFlag == "" {
			log.Fatal("Для --get нужны --id")
		}
		raw, err := apiClient.GetBin(*idFlag)
		if err != nil {
			log.Fatalf("Ошибка при получении: %v", err)
		}
		var formatted bytes.Buffer
		if err := json.Indent(&formatted, raw, "", "  "); err != nil {
			log.Fatalf("Ошибка форматирования JSON: %v", err)
		}
		fmt.Println("Содержимое бина:")
		fmt.Println(formatted.String())
	case *deleteCmd:
		if *idFlag == "" {
			log.Fatal("Для --delete нужны --id")
		}
		if err := apiClient.DeleteBin(*idFlag); err != nil {
			log.Fatalf("Ошибка удаления: %v", err)
		}
		fmt.Println("Бин успешно удалён.")
		if err := store.RemoveFromIndex(*idFlag); err != nil {
			log.Printf("Не удалось удалить бин из индекса: %v\n", err)
		}
	case *listCmd:
		index, err := store.GetIndex()
		if err != nil {
			log.Fatalf("Не удалось получить список бинoв: %v", err)
		}
		if len(index) == 0 {
			fmt.Println("Бины не найдены.")
		} else {
			fmt.Println("Список бинов:")
			for name, id := range index {
				fmt.Printf("- %s: %s\n", name, id)
			}
		}
	default:
		fmt.Println("Неизвестная или неуказанная команда. Используй --create, --update, --get, --delete или --list")
	}
}
