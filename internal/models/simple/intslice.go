package simple

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Custom IntSlice type for handling "1+2+3+4+5" storage
type IntSlice []int

// Scan converts database value to []int
func (s *IntSlice) Scan(value interface{}) error {

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to convert database value to bytes")
	}

	str := string(bytes)

	*s = *ParseIntSlice(str)
	return nil
}

// Value converts []int to "1+2+3+4+5" for storage
func (s IntSlice) Value() (driver.Value, error) {
	return s.String(), nil
}

// String converts []int to "1+2+3+4+5"
func (s IntSlice) String() string {
	var strParts []string
	for _, num := range s {
		strParts = append(strParts, strconv.Itoa(num))
	}
	return strings.Join(strParts, "+")
}

func ParseIntSlice(str string) *IntSlice {
	parts := strings.Split(str, "+")
	var result IntSlice
	for _, part := range parts {
		if num, err := strconv.Atoi(part); err == nil {
			result = append(result, num)
		}
	}

	return &result
}
