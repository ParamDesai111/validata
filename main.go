package main

import (
	"validata/utils"
	"validata/validator"
	"fmt"
)

func main() {
	data, err := utils.ReadCSV("data/sample.csv", true)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	missingCounts := validator.DetectMissingValues(data)
	fmt.Println("Missing Values Count:", missingCounts)

	imputedData := validator.ImputeMissingValues(data, "mean")
	fmt.Println("Imputed Data:", imputedData)
}
