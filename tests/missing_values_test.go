package tests

import (
	"fmt"
	"math"
	"testing"
	"validata/utils"
	"validata/validator"

	"github.com/go-gota/gota/dataframe"
)

// TestReadCSV checks if the CSV file is read correctly with headers
func TestReadCSV_withHeader(t *testing.T) {
	data, err := utils.ReadCSV("sample_data/test.csv", true)
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}
	fmt.Println("CSV Data:", data)

	// Check that the expected number of rows is returned (excluding header)
	expectedRows := 5
	if len(data) != expectedRows {
		t.Errorf("Expected %d rows, but got %d", expectedRows, len(data))
	}
}

// TestReadCSV checks if the CSV file is read correctly without headers
func TestReadCSV_withoutHeader(t *testing.T) {
	data, err := utils.ReadCSV("sample_data/test.csv", false)
	if err != nil {
		t.Fatalf("Failed to read CSV: %v", err)
	}
	fmt.Println("CSV Data:", data)

	// Check that the expected number of rows is returned (including the header)
	expectedRows := 6
	if len(data) != expectedRows {
		t.Errorf("Expected %d rows, but got %d", expectedRows, len(data))
	}
}

// TestReadJSON checks if the JSON file is read correctly with headers
func TestReadJSON_Headers(t *testing.T) {
	data, err := utils.ReadJSON("sample_data/test.json", false)
	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}
	fmt.Println("JSON Data:", data)

	// Check that the expected number of rows is returned
	expectedRows := 5
	if len(data) != expectedRows {
		t.Errorf("Expected %d rows, but got %d", expectedRows, len(data))
	}
}

func TestReadJSON_noHeaders(t *testing.T) {
	data, err := utils.ReadJSON("sample_data/test.json", true)
	if err != nil {
		t.Fatalf("Failed to read JSON: %v", err)
	}
	fmt.Println("JSON Data:", data)

	// Check that the expected number of rows is returned
	expectedRows := 6
	if len(data) != expectedRows {
		t.Errorf("Expected %d rows, but got %d", expectedRows, len(data))
	}
}

// TestDetectMissingValues validates missing value detection for CSV data
func TestDetectMissingValues(t *testing.T) {
	data, _ := utils.ReadCSV("sample_data/test.csv", true)
	missing := validator.DetectMissingValues(data)

	fmt.Println("Missing Values Count:", missing)

	expected := []int{0, 1, 1, 1} // Expected missing values per column
	for i, v := range missing {
		if v != expected[i] {
			t.Errorf("Expected %d missing values in column %d, but got %d", expected[i], i, v)
		}
	}
}

// TestImputeMissingValues verifies missing value imputation for CSV data
func TestImputeMissingValues(t *testing.T) {
	data, _ := utils.ReadCSV("sample_data/test.csv", false)

	imputedData := validator.ImputeMissingValues(data, "mean")
	fmt.Println("Imputed Data:", imputedData)

	// Check if missing values have been filled
	for row := 0; row < len(imputedData); row++ {
		for col := 0; col < len(imputedData[row]); col++ {
			if imputedData[row][col] == "" {
				t.Errorf("Expected missing value at row %d, col %d to be imputed, but it's still empty", row, col)
			}
		}
	}
}

// TestDetectMissingValuesJSON validates missing value detection for JSON data
func TestDetectMissingValuesJSON(t *testing.T) {
	data, _ := utils.ReadJSON("sample_data/test.json", false)
	missing := validator.DetectMissingValues(data)

	fmt.Println("Missing Values Count (JSON):", missing)

	expected := []int{0, 1, 1, 1}
	for i, v := range missing {
		if v != expected[i] {
			t.Errorf("Expected %d missing values in column %d, but got %d", expected[i], i, v)
		}
	}
}


// TestDetectMissingValuesDF validates missing values detection in DataFrames
func TestDetectMissingValuesDF(t *testing.T) {
	df := dataframe.LoadRecords(
		[][]string{
			{"ID", "Name", "Age", "Salary"},
			{"1", "Alice", "25", "50000"},
			{"2", "Bob", "", "60000"},
			{"3", "", "30", "70000"},
			{"4", "David", "40", ""},
			{"5", "Eve", "35", "80000"},
		},
	)

	missingCounts := validator.DetectMissingValuesDF(df)
	fmt.Println("Missing Values Count (DataFrame):", missingCounts)

	expected := map[string]int{"ID": 0, "Name": 1, "Age": 1, "Salary": 1}
	for col, count := range missingCounts {
		if count != expected[col] {
			t.Errorf("Expected %d missing values in column %s, but got %d", expected[col], col, count)
		}
	}
}


// TestImputeMissingValuesDF verifies missing value imputation in DataFrames
func TestImputeMissingValuesDF(t *testing.T) {
	df := dataframe.LoadRecords(
		[][]string{
			{"ID", "Name", "Age", "Salary"},
			{"1", "Alice", "25", "50000"},
			{"2", "Bob", "", "60000"},
			{"3", "", "30", "70000"},
			{"4", "David", "40", ""},
			{"5", "Eve", "35", "80000"},
		},
	)

	imputedDF := validator.ImputeMissingValuesDF(df, "mean")
	fmt.Println("Imputed DataFrame:", imputedDF)

	// Ensure missing values are correctly replaced
	for _, col := range df.Names() {
		colData := df.Col(col)
		for i := 0; i < colData.Len(); i++ {
			val := colData.Elem(i).Float()
			if math.IsNaN(val) {
				t.Errorf("Expected missing value in column %s to be imputed, but it's still NaN", col)
			}
		}
	}
}
