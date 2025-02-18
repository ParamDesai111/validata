package validator

import (
	"fmt"
	"math"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// DetectMissingValues checks for missing values in a dataset
func DetectMissingValues(data [][]string) []int {
	if len(data) == 0 {
		return nil // Prevent index out of range errors
	}

	missingCounts := make([]int, len(data[0]))

	for col := 0; col < len(data[0]); col++ {
		missingCount := 0

		for row := 0; row < len(data); row++ {
			val := data[row][col]

			// Check if value is missing (empty, "NULL", or "null")
			if val == "" || val == "NULL" || val == "null" {
				missingCount++
			}
		}
		missingCounts[col] = missingCount
	}

	return missingCounts
}

// ImputeMissingValues fills missing values with a given strategy (mean, median)
func ImputeMissingValues(data [][]string, strategy string) [][]string {
	columnData := make([][]float64, len(data[0]))

	for col := 0; col < len(data[0]); col++ {
		var colValues []float64

		for row := 0; row < len(data); row++ {
			if data[row][col] == "" || data[row][col] == "NULL" {
				continue
			}

			var val float64
			_, err := fmt.Sscanf(data[row][col], "%f", &val)
			if err == nil {
				colValues = append(colValues, val)
			}
		}
		columnData[col] = colValues
	}

	for col := 0; col < len(data[0]); col++ {
		var replacementValue float64

		switch strategy {
		case "mean":
			replacementValue = Mean(columnData[col])
		case "median":
			replacementValue = Median(columnData[col])
		default:
			fmt.Println("Unsupported strategy, defaulting to mean")
			replacementValue = Mean(columnData[col])
		}

		for row := 0; row < len(data); row++ {
			if data[row][col] == "" || data[row][col] == "NULL" {
				data[row][col] = fmt.Sprintf("%.2f", replacementValue)
			}
		}
	}

	return data
}

// Mean calculates the mean of a slice of numbers
func Mean(nums []float64) float64 {
	if len(nums) == 0 {
		return math.NaN()
	}
	sum := 0.0
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

// Median calculates the median of a slice of numbers
func Median(nums []float64) float64 {
	if len(nums) == 0 {
		return math.NaN()
	}
	n := len(nums)
	mid := n / 2
	if n%2 == 0 {
		return (nums[mid-1] + nums[mid]) / 2
	}
	return nums[mid]
}

// DetectMissingValuesDF checks for missing values in a DataFrame
func DetectMissingValuesDF(df dataframe.DataFrame) map[string]int {
	missingCounts := make(map[string]int)

	for _, colName := range df.Names() {
		missingCount := 0
		col := df.Col(colName)

		for i := 0; i < col.Len(); i++ {
			// Check for empty strings (""), "NULL", or "null" in any column
			val := col.Elem(i).String()
			if val == "" || val == "NULL" || val == "null" {
				missingCount++
			}
		}
		missingCounts[colName] = missingCount
	}

	return missingCounts
}


// ImputeMissingValuesDF fills missing values with a given strategy (mean, median)
func ImputeMissingValuesDF(df dataframe.DataFrame, strategy string) dataframe.DataFrame {
	imputedCols := make([]series.Series, 0) // Holds updated columns

	for _, colName := range df.Names() {
		col := df.Col(colName)

		// Handle numeric columns
		if col.Type() == series.Float {
			var numericValues []float64
			for i := 0; i < col.Len(); i++ {
				valStr := col.Elem(i).String()
				if valStr != "" { // Ignore empty strings
					val := col.Elem(i).Float()
					numericValues = append(numericValues, val)
				}
			}

			// Compute replacement value
			var replacementValue float64
			if len(numericValues) > 0 {
				switch strategy {
				case "mean":
					replacementValue = Mean(numericValues)
				case "median":
					replacementValue = Median(numericValues)
				default:
					fmt.Println("Unsupported strategy, defaulting to mean")
					replacementValue = Mean(numericValues)
				}
			} else {
				replacementValue = 0.0 // Default replacement if no valid values exist
			}

			// Replace empty values
			newCol := make([]float64, col.Len())
			for i := 0; i < col.Len(); i++ {
				valStr := col.Elem(i).String()
				if valStr == "" { // Check for empty string instead of NaN
					newCol[i] = replacementValue // Replace "" with computed value
				} else {
					newCol[i] = col.Elem(i).Float()
				}
			}

			// Append the modified column
			imputedCols = append(imputedCols, series.New(newCol, series.Float, colName))
		} else if col.Type() == series.String {
			// Handle string columns
			newCol := make([]string, col.Len())
			for i := 0; i < col.Len(); i++ {
				val := col.Elem(i).String()
				if val == "" || val == "NULL" || val == "null" {
					newCol[i] = "Unknown" // Replace empty strings with "Unknown"
				} else {
					newCol[i] = val
				}
			}
			imputedCols = append(imputedCols, series.New(newCol, series.String, colName))
		} else {
			// Retain original column if it's neither numeric nor string
			imputedCols = append(imputedCols, col)
		}
	}

	// ✅ Ensure all columns are retained in the DataFrame
	return dataframe.New(imputedCols...)
}
