package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const (
	opAvg = "AVG"
	opSum = "SUM"
	opMed = "MED"
)

func main() {
	var operationType string
	var numbersString string
	var numbersArray []int64

	for {
		fmt.Println("Input operation type (AVG/SUM/MED)")
		_, err := fmt.Scan(&operationType)
		operationType = strings.ToUpper(operationType)
		if err == nil && (operationType == opAvg || operationType == opSum || operationType == opMed) {
			break
		}
		fmt.Println("Wrong operation type, try again")
	}

	for {
		fmt.Println("Enter numbers separated by commas")
		_, err := fmt.Scanln(&numbersString)
		if err != nil {
			fmt.Println("Error reading input, please try again:", err)
			continue
		} else {
			break
		}
	}
	numbersArray = convertStringToArray(numbersString)
	if len(numbersArray) == 0 {
		fmt.Println("No numbers provided")
	} else {
		switch operationType {
		case opAvg:

			fmt.Println("AVG = ", getAvg(numbersArray))
		case opSum:
			fmt.Println("SUM = ", getSum(numbersArray))
		case opMed:
			fmt.Println("MED = ", getMed(numbersArray))
		}
	}
}

func getAvg(numbersArray []int64) float64 {
	sum := getSum(numbersArray)
	average := float64(sum) / float64(len(numbersArray))
	return average

}

func getSum(numbersArray []int64) int64 {
	var sum int64
	for _, num := range numbersArray {
		sum += num
	}
	return sum
}

func getMed(numbersArray []int64) float64 {
	sortedArray := make([]int64, len(numbersArray))
	copy(sortedArray, numbersArray)

	sort.Slice(sortedArray, func(i, j int) bool {
		return sortedArray[i] < sortedArray[j]
	})

	length := len(sortedArray)

	if length%2 == 0 {
		middle1 := sortedArray[length/2-1]
		middle2 := sortedArray[length/2]
		return float64(middle1+middle2) / 2.0
	}

	return float64(sortedArray[length/2])

}

func convertStringToArray(numbersString string) []int64 {
	strs := strings.Split(numbersString, ",")
	var result []int64
	for _, str := range strs {
		num, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
		if err == nil {
			result = append(result, num)
		}
	}
	return result
}
