package main

import "fmt"

func main() {
	const USDtoEUR = 0.93
	const USDtoRUB = 75
	const EURtoRUB = USDtoRUB / USDtoEUR

	fmt.Println("EUR to RUB exchange rate:", EURtoRUB)
}
