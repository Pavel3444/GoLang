package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

type CurrencyMap map[string]map[string]float64

var currencyTree = &CurrencyMap{
	"USD": {"EUR": 0.93, "RUB": 75.0},
	"EUR": {"USD": 1.11, "RUB": 86.0},
	"RUB": {"USD": 0.0133, "EUR": 0.01239},
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	currency := askCurrency(currencyTree)
	amount, errAmount := askAmount()
	if errAmount != nil {
		return errAmount
	}
	targetCurrency, errTargetCurrency := askTargetCurrency(currency)
	if errTargetCurrency != nil {
		return errTargetCurrency
	}
	res, errConvert := convert(currencyTree, currency, targetCurrency, amount)
	if errConvert != nil {
		return errConvert
	}
	fmt.Printf("Вы конвертируете: %.2f %s в %.2f %s\n", amount, currency, res, targetCurrency)
	return nil
}

func askCurrency(tree *CurrencyMap) string {
	var currency string
	curreces := make([]string, 0, len(*tree))
	for k := range *tree {
		curreces = append(curreces, k)
	}
	for {
		fmt.Println("Введите код валюты для конвертирования (USD, EUR, RUB): ")
		_, err := fmt.Scan(&currency)
		if err == nil && contains(curreces, currency) {
			break
		}
		if err != nil {
			clearBuffer()
		}
		fmt.Println("Вы ввели неверный код валюты, заново введите код из предложенных, соблюдая регистр")
	}
	return currency
}

func askAmount() (float64, error) {
	var amount float64
	fmt.Print("Введите сумму конвертирования: ")
	_, err := fmt.Scan(&amount)
	if err != nil {
		clearBuffer()
		return 0, fmt.Errorf("введено не число")
	}
	if amount <= 0 {
		return 0, fmt.Errorf("сумма должна быть больше нуля")
	}
	return amount, nil
}

func askTargetCurrency(currency string) (string, error) {
	var targetCurrency string
	targets, err := getTargetCurrency(currencyTree, currency)
	if err != nil {
		return "", err
	} else {
		list := strings.Join(targets, ", ")
		for {
			fmt.Println("Введите код целевой валюты для конвертирования: ", list)
			_, err := fmt.Scan(&targetCurrency)
			if err == nil && contains(targets, targetCurrency) {
				break
			}
			if err != nil {
				clearBuffer()
			}
			fmt.Println("Неверная валюта, введите одну из следующих, соблюдая регистр: ", list)
		}
		return targetCurrency, nil
	}
}

func convert(tree *CurrencyMap, currency, targetCurrency string, amount float64) (float64, error) {
	ratesForCurrency, ok := (*tree)[currency]
	if !ok {
		return 0, fmt.Errorf("исходная валюта %q не найдена", currency)
	}
	rate, ok := ratesForCurrency[targetCurrency]
	if !ok {
		return 0, fmt.Errorf("целевая валюта %q не найдена", targetCurrency)
	}
	res := amount * rate
	if math.IsNaN(res) || math.IsInf(res, 0) {
		return 0, fmt.Errorf("ошибка конвертации: результат недействителен")
	}

	return math.Round(res*100) / 100, nil
}

func getTargetCurrency(tree *CurrencyMap, currency string) ([]string, error) {
	targets, exists := (*tree)[currency]
	if !exists {
		return nil, fmt.Errorf("валюта %s не найдена", currency)
	}
	keys := make([]string, 0, len(targets))
	for k := range targets {
		keys = append(keys, k)
	}
	return keys, nil
}

func clearBuffer() {
	var discard string
	fmt.Scanln(&discard)
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}
