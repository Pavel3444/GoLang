package main

import (
	"fmt"
	"time"

	"bin/bins"
	"bin/storage"
)

const fileName = "bins.json"

func main() {
	list := bins.NewBinList()

	list.AddBin(false, "Firstbin")
	list.AddBin(true, "Secondbin")

	data, err := list.ToBytes()
	if err != nil {
		fmt.Printf("Error serializing bins: %v\n", err)
		return
	}

	if err := storage.SaveBins(data, fileName); err != nil {
		fmt.Printf("Error saving bins: %v\n", err)
		return
	}
	fmt.Printf("Bins successfully saved to %q\n", fileName)

	loaded := bins.NewBinList()
	if _, err := storage.ReadBins(fileName, loaded); err != nil {
		fmt.Printf("Error reading bins: %v\n", err)
		return
	}

	fmt.Println("Loaded bins:")
	for _, b := range loaded.Bins {
		fmt.Printf("- ID: %s, Name: %s, Private: %t, Created: %s\n",
			b.Id, b.Name, b.Private, b.CreateTime.Format(time.RFC3339))
	}
}
