package gift

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type GiftList []*Gift

// Scan implements the sql.Scanner interface
func (gl *GiftList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	str, ok := value.(string)
	if !ok {
		bytes, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("failed to unmarshal Avatar value: %v", value)
		}
		str = string(bytes)
	}

	gifts := []*Gift{}
	for _, giftStr := range strings.Split(str, "#") {
		if giftStr == "" {
			continue
		}
		gift, err := NewGiftFromString(giftStr)
		if err != nil {
			return err
		}
		gifts = append(gifts, gift)
	}

	var giftList GiftList = gifts

	*gl = giftList
	return nil
}

// Value implements the driver.Valuer interface
func (gl GiftList) Value() (driver.Value, error) {
	return gl.String(), nil
}

// String turns the gift list to string
func (gl GiftList) String() string {
	giftsStr := []string{}
	for _, gift := range gl {
		giftsStr = append(giftsStr, gift.String())
	}

	return strings.Join(giftsStr, "#")
}
