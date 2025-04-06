package main

import "fmt"

func main() {
	const USDtoEUR = 0.93
	const USDtoRUB = 75
	const EURtoRUB = USDtoRUB / USDtoEUR

	currency, targetCurrency, amount := userInputData()

	fmt.Println("your data:", currency, targetCurrency, amount)

}

func userInputData() (string, string, float64) {
	var currency string
	var targetCurrency string
	var amount float64

	fmt.Println("Введите код валюты для конвертирования: ")
	fmt.Scan(&currency)
	fmt.Println("Введите код целевой валюты для конвертирования: ")
	fmt.Scan(&targetCurrency)
	fmt.Println("Введите сумму конвертирования: ")
	fmt.Scan(&amount)

	return currency, targetCurrency, amount
}

func convert(currency, targetCurrency string, amount float64) float64 {
	return amount
}
