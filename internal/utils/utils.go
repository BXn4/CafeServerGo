/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package utils

import (
	"fmt"
	"strconv"
	"time"
)

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
	now := time.Now().UTC()

	for _, arg := range args {
		if val, err := strconv.Atoi(arg); err == nil {
			result = append(result, val)
			continue
		}

		if t, err := time.Parse(time.RFC3339, arg); err == nil {
			diff := int(t.Sub(now).Seconds())
			result = append(result, diff)
			continue
		}

		return nil, fmt.Errorf("Cannot parse %q as int or datetime", arg)
	}

	return result, nil
}
