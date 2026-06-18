package simple

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// IntMap stores map[int]int as a string "1+100#2+200"
type IntMap map[int]int

// Scan converts a DB value into IntMap
func (m *IntMap) Scan(value interface{}) error {
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case nil:
		*m = make(IntMap)
		return nil
	default:
		return fmt.Errorf("failed to convert database value to bytes: unsupported type %T", value)
	}

	*m = *ParseIntMap(str)
	return nil
}

// Value converts IntMap to a database-storable string
func (m IntMap) Value() (driver.Value, error) {
	return m.String(), nil
}

func (m IntMap) String() string {
	var entries []string
	for key, value := range m {
		entries = append(entries, fmt.Sprintf("%d+%d", key, value))
	}
	return strings.Join(entries, "#")
}

func ParseIntMap(str string) *IntMap {
	m := IntMap{}
	if str == "" {
		return &m
	}

	entries := strings.Split(str, "#")
	for _, entry := range entries {
		parts := strings.Split(entry, "+")
		if len(parts) != 2 {
			continue
		}

		key, err1 := strconv.Atoi(parts[0])
		value, err2 := strconv.Atoi(parts[1])
		if err1 == nil && err2 == nil {
			m[key] = value
		}
	}

	return &m
}
