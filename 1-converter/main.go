package main

import (
	"fmt"
)

const USDtoEUR = 0.93
const USDtoRUB = 75
const EURtoRUB = USDtoRUB / USDtoEUR

func main() {
	currency := askCurrency()
	amount := askAmount()
	targetCurrency := askTargetCurrency(currency)
	res := convert(currency, targetCurrency, amount)

	fmt.Println("Вы конвертируете :", amount, currency, "в ", res, targetCurrency)
}

func askCurrency() string {
	var currency string
	for {
		fmt.Println("Введите код валюты для конвертирования (USD, EUR, RUB): ")
		_, err := fmt.Scan(&currency)
		if err == nil && (currency == "USD" || currency == "EUR" || currency == "RUB") {
			break
		}
		if err != nil {
			clearBuffer()
		}
		fmt.Println("Вы ввели неверный код валюты, заново введите код из предложенных, соблюдая регистр")
	}
	return currency
}

func askAmount() float64 {
	var amount float64
	for {
		fmt.Println("Введите сумму конвертирования:")
		_, err := fmt.Scan(&amount)
		if err == nil && amount > 0 {
			break
		}
		if err != nil {
			clearBuffer()
		}
		fmt.Println("Вы ввели неверную сумму, введите заново")
	}
	return amount
}

func askTargetCurrency(currency string) string {
	var targetCurrency string
	currencyOption1, currencyOption2 := getTargetCurrency(currency)
	for {
		fmt.Println("Введите код целевой валюты для конвертирования: ", currencyOption1, currencyOption2)
		_, err := fmt.Scan(&targetCurrency)
		if err == nil && (targetCurrency == currencyOption1 || targetCurrency == currencyOption2) {
			break
		}
		if err != nil {
			clearBuffer()
		}
		fmt.Println("Неверная валюта, введите одну из следующих, соблюдая регистр: ", currencyOption1, currencyOption2)
	}
	return targetCurrency
}

func convert(currency, targetCurrency string, amount float64) float64 {
	switch {
	case currency == "USD" && targetCurrency == "EUR":
		return amount * USDtoEUR
	case currency == "USD" && targetCurrency == "RUB":
		return amount * USDtoRUB
	case currency == "EUR" && targetCurrency == "USD":
		return amount / USDtoEUR
	case currency == "EUR" && targetCurrency == "RUB":
		return amount * EURtoRUB
	case currency == "RUB" && targetCurrency == "USD":
		return amount / USDtoRUB
	case currency == "RUB" && targetCurrency == "EUR":
		return amount / EURtoRUB
	default:
		fmt.Println("Некорректная операция конвертации")
		return 0
	}
}

func getTargetCurrency(currency string) (string, string) {
	switch currency {
	case "USD":
		return "EUR", "RUB"
	case "EUR":
		return "USD", "RUB"
	default:
		return "USD", "EUR"
	}
}

func clearBuffer() {
	var discard string
	fmt.Scanln(&discard)
}
