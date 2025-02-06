package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// ReadCSV reads a CSV file and converts it into a 2D string array
func ReadCSV(filePath string, hasHeader bool) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// If the file has a header, remove the first row
	if hasHeader && len(rows) > 0 {
		rows = rows[1:]
	}

	return rows, nil
}

// ReadJSON reads a JSON file and converts it into a 2D string array
func ReadJSON(filePath string, hasHeader bool) ([][]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var jsonData []map[string]interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		return nil, err
	}

	if len(jsonData) == 0 {
		return nil, fmt.Errorf("empty JSON file")
	}

	// Extract column headers
	var headers []string
	for key := range jsonData[0] {
		headers = append(headers, key)
	}

	// Convert JSON to 2D string format
	var data [][]string

	// If hasHeader is true, the first row contains headers
	if hasHeader {
		data = append(data, headers)
	}

	for _, obj := range jsonData {
		var row []string
		for _, key := range headers {
			row = append(row, fmt.Sprintf("%v", obj[key]))
		}
		data = append(data, row)
	}

	return data, nil
}