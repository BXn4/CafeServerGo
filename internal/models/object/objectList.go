package object

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type ObjectList []*Object

// Scan implements the sql.Scanner interface
func (ol *ObjectList) Scan(value interface{}) error {

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("\nFailed to convert database value to bytes: %v", value)
	}

	str := string(bytes)
	if str == "" {
		println("Theres no objects found :(")
		objectList := ObjectList{}
		*ol = objectList
		return nil
	}

	objectList := ParseObjectList(str)
	if objectList == nil {
		return fmt.Errorf("\nFailed to parse database value to ObjectList")
	}

	*ol = *objectList
	return nil
}

// Value implements the driver.Value interface
func (ol ObjectList) Value() (driver.Value, error) {
	return ol.String(), nil
}

// String turns the ObjectList to string
func (ol ObjectList) String() string {
	objsStr := []string{}
	for _, obj := range ol {
		objsStr = append(objsStr, obj.String())
	}

	return strings.Join(objsStr, "#")
}

func (ol ObjectList) StringForDB() string {
	objsStr := []string{}
	for _, obj := range ol {
		objsStr = append(objsStr, obj.StringForDB())
	}

	return strings.Join(objsStr, "#")
}

func ParseObjectList(str string) *ObjectList {
	parts := strings.Split(str, "#")
	var result ObjectList
	for _, objStr := range parts {
		obj, err := NewObjectFromString(objStr)
		if err != nil {
			return nil
		}
		result = append(result, obj)
	}

	return &result
}
