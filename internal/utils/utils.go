package utils

import "strconv"

// This is a ternary operator
//
//	utils.If(true, "foo", "bar") // returns "foo"
//	utils.If(false, "foo", "bar") // returns "bar"
func If[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

func MultiAtoi(args ...string) ([]int, error) {
	result := []int{}
	for _, arg := range args {
		val, err := strconv.Atoi(arg)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}
