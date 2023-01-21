package service

import (
	"errors"
	"strings"
)

func getQuestionMarks(values []any) (string, error) {
	if len(values) == 0 {
		return "", errors.New("empty array of values")
	}
	length := len(values)
	slice := make([]string, length)
	for i := range slice {
		slice[i] = "?"
	}
	return strings.Join(slice, ", "), nil
}

func GetSliceKeys(m map[int]interface{}) []int {
	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
func ArrayValues(elements map[interface{}]interface{}) []interface{} {
	i, values := 0, make([]interface{}, len(elements))
	for _, val := range elements {
		values[i] = val
		i++
	}
	return values
}