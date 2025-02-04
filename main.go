package main

import (
	"fmt"
	"validata/validator"
)

func main() {
	fmt.Println("Data Quality Toolkit Running...")

	// Example Data
	data := [][]string{
		{"1", "2", ""},
		{"4", "", "6"},
		{"7", "8", "9"},
	}

	missing := validator.DetectMissingValues(data)
	fmt.Println("Missing Values Count:", missing)
}
