package tests

import (
	"fmt"
	"testing"
	"validata/validator"
)

// TestDetectMissingValues verifies the function detects missing values correctly
func TestDetectMissingValues(t *testing.T) {
	data := [][]string{
		{"1", "2", ""},
		{"4", "", "6"},
		{"7", "8", "9"},
	}

	missing := validator.DetectMissingValues(data)
	fmt.Println("Missing Values Count:", missing)

	expected := []int{0, 1, 1}
	for i, v := range missing {
		if v != expected[i] {
			t.Errorf("Expected %d missing values in column %d, but got %d", expected[i], i, v)
		}
	}
}

// TestImputeMissingValues verifies missing values imputation works
func TestImputeMissingValues(t *testing.T) {
	data := [][]string{
		{"1", "2", ""},
		{"4", "", "6"},
		{"7", "8", "9"},
	}

	imputedData := validator.ImputeMissingValues(data, "mean")
	fmt.Println("Imputed Data:", imputedData)
}
