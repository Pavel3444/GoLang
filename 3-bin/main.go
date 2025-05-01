package main

import (
	"bin/bins"
	"bin/file"
	"bin/storage"
	"log"
)

func main() {
	fileClient := file.NewFile()
	store := storage.NewStorage(fileClient)

	binRepo := bins.NewBinRepository(store, "bins.json")

	b, err := binRepo.Add(false, "first bin test")
	if err != nil {
		log.Fatalf("cannot add bin: %v", err)
	}
	log.Printf("Created bin: %+v\n", b)

	all, err := binRepo.List()
	if err != nil {
		log.Fatalf("cannot list bins: %v", err)
	}
	log.Printf("All bins: %+v\n", all)
}
