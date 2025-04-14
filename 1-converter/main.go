package main

import (
	"fmt"
	"math"
)

var currencyTree = map[string]map[string]float64{
	"USD": {
		"EUR": 0.93,
		"RUB": 75.0,
	},

	"EUR": {
		"USD": 1.11,
		"RUB": 86.0,
	},
	"RUB": {
		"USD": 0.0133,
		"EUR": 0.01239,
	},
}

func main() {
	currency := askCurrency()
	amount := askAmount()
	targetCurrency := askTargetCurrency(currency)
	res := convert(currency, targetCurrency, amount)

	fmt.Printf("Вы конвертируете: %.2f %s в %.2f %s\n", amount, currency, res, targetCurrency)

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
	res := amount * currencyTree[currency][targetCurrency]
	if math.IsNaN(res) {
		fmt.Println("Ошибка конвертации: результат недействителен")
		return 0
	}

	return math.Round(res*100) / 100
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
