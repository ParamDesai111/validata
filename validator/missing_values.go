package validator

import (
	"fmt"
	"math"
)

// DetectMissingValues checks for missing values in a dataset
func DetectMissingValues(data [][]string) []int {
	missingCounts := make([]int, len(data[0]))

	// Iterate over each column
	for col := 0; col < len(data[0]); col++ {
		missingCount := 0

		// Iterate over each row
		for row := 0; row < len(data); row++ {
			if data[row][col] == "" || data[row][col] == "NULL" {
				missingCount++
			}
		}
		missingCounts[col] = missingCount
	}

	return missingCounts
}

// ImputeMissingValues fills missing values with a given strategy (mean, median, mode)
func ImputeMissingValues(data [][]string, strategy string) [][]string {
	columnData := make([][]float64, len(data[0]))

	// Convert data into float format where possible
	for col := 0; col < len(data[0]); col++ {
		var colValues []float64

		for row := 0; row < len(data); row++ {
			if data[row][col] == "" || data[row][col] == "NULL" {
				continue // Skip missing values
			}

			var val float64
			_, err := fmt.Sscanf(data[row][col], "%f", &val)
			if err == nil {
				colValues = append(colValues, val)
			}
		}
		columnData[col] = colValues
	}

	// Apply chosen imputation strategy
	for col := 0; col < len(data[0]); col++ {
		var replacementValue float64

		switch strategy {
		case "mean":
			replacementValue = Mean(columnData[col])
		case "median":
			replacementValue = Median(columnData[col])
		default:
			fmt.Println("Unsupported imputation strategy, defaulting to mean")
			replacementValue = Mean(columnData[col])
		}

		// Replace missing values with computed value
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

	// Sort values
	n := len(nums)
	mid := n / 2
	if n%2 == 0 {
		return (nums[mid-1] + nums[mid]) / 2
	}
	return nums[mid]
}
