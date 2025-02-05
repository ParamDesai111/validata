package validator

import (
	"fmt"
	"math"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// DetectMissingValues checks for missing values in a dataset
func DetectMissingValues(data [][]string) []int {
	missingCounts := make([]int, len(data[0]))

	for col := 0; col < len(data[0]); col++ {
		missingCount := 0

		for row := 0; row < len(data); row++ {
			if data[row][col] == "" || data[row][col] == "NULL" {
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
			val := col.Elem(i).String()
			if val == "" || val == "NULL" {
				missingCount++
			}
		}
		missingCounts[colName] = missingCount
	}

	return missingCounts
}

// ImputeMissingValuesDF fills missing values with a given strategy (mean, median)
func ImputeMissingValuesDF(df dataframe.DataFrame, strategy string) dataframe.DataFrame {
	for _, colName := range df.Names() {
		col := df.Col(colName)

		// Collect numeric values
		var numericValues []float64
		for i := 0; i < col.Len(); i++ {
			val := col.Elem(i).Float()
			if !math.IsNaN(val) {
				numericValues = append(numericValues, val)
			}
		}

		// Compute replacement value
		var replacementValue float64
		switch strategy {
		case "mean":
			replacementValue = Mean(numericValues)
		case "median":
			replacementValue = Median(numericValues)
		default:
			fmt.Println("Unsupported strategy, defaulting to mean")
			replacementValue = Mean(numericValues)
		}

		// Create a new column with imputed values
		newCol := make([]float64, col.Len())
		for i := 0; i < col.Len(); i++ {
			val := col.Elem(i).Float()
			if math.IsNaN(val) {
				newCol[i] = replacementValue // Replace NaN with computed value
			} else {
				newCol[i] = val
			}
		}

		// Update the DataFrame with the new column
		df = df.Mutate(series.Floats(newCol)).Rename(colName, fmt.Sprintf("%s_imputed", colName))
	}

	return df
}