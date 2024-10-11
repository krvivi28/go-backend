package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ConvertPostgresArrayTo2D(arrayStr string) ([][]int, error) {
	arrayStr = strings.Trim(arrayStr, "{}")

	rowStrings := strings.Split(arrayStr, "},{")
	var result [][]int

	for _, rowStr := range rowStrings {
		rowStr = strings.Trim(rowStr, "{}")

		elements := strings.Split(rowStr, ",")
		var row []int

		for _, elem := range elements {
			val, err := strconv.Atoi(elem)
			if err != nil {
				return nil, err
			}
			row = append(row, val)
		}
		result = append(result, row)
	}

	return result, nil
}

func ConvertToPostgresArray(data [][]int) string {
	var rows []string
	for _, row := range data {
		strRow := make([]string, len(row))
		for i, val := range row {
			strRow[i] = fmt.Sprintf("%d", val)
		}
		rows = append(rows, fmt.Sprintf("{%s}", strings.Join(strRow, ",")))
	}
	return fmt.Sprintf("{%s}", strings.Join(rows, ","))
}
