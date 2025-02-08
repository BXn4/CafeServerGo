package simple

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

type Position struct {
	X int
	Y int
}

func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

// Value implements the driver.Valuer interface for database storage
func (p Position) Value() (driver.Value, error) {
	// Format: "x+y"
	return fmt.Sprintf("%d+%d", p.X, p.Y), nil
}

// Scan implements the sql.Scanner interface for database retrieval
func (p *Position) Scan(value interface{}) error {
	if value == nil {
		*p = Position{} // Set zero value
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to scan Position: unexpected type %T", value)
		}
		str = string(bytes)
	}

	// Parse the string format "x+y"
	parts := strings.Split(str, "+")
	if len(parts) != 2 {
		return fmt.Errorf("invalid Position format: %s", str)
	}

	var x, y int
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse x coordinate: %v", err)
	}

	y, err = strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("failed to parse y coordinate: %v", err)
	}

	*p = Position{
		X: x,
		Y: y,
	}

	return nil
}
