package waiter

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type WaiterList []*Waiter

// Scan converts a database value (string) to IntMatrix
func (wl *WaiterList) Scan(value interface{}) error {
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case nil:
		waiterList := &WaiterList{}
		*wl = *waiterList
		return nil
	default:
		return fmt.Errorf("failed to convert database value to bytes")
	}

	waiterList := ParseWaiterList(str)
	if waiterList == nil {
		return fmt.Errorf("\nFailed to parse database value to WaiterList")
	}

	*wl = *waiterList
	return nil
}

func (wl WaiterList) Value() (driver.Value, error) {
	if len(wl) == 0 {
		return "", nil
	}

	return wl.String(), nil
}

func (wl WaiterList) String() string {
	var result []string
	for _, w := range wl {
		result = append(result, w.String())
	}

	return strings.Join(result, "%")
}

func ParseWaiterList(str string) *WaiterList {

	parts := strings.Split(str, "%")
	var result WaiterList
	for _, wStr := range parts {
		w := NewWaiterFromString(wStr)
		if w == nil {
			continue
		}
		result = append(result, w)
	}

	return &result
}
