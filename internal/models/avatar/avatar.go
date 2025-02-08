package avatar

import (
	"cafego/internal/utils"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Parts int

const (
	TOP  Parts = 1000
	SKIN       = 1020
	HAIR       = 1040
	LEGS       = 1050
	HAT        = 1060
	FACE       = 1080
)

type AvatarGender int

const (
	Girl AvatarGender = 1
	Boy               = 2
)

type Avatar struct {
	ID        int          `gorm:"type:int;not null"`
	Name      string       `gorm:"type:string;not null"`
	Gender    AvatarGender `gorm:"type:int;not null"`
	SkinColor int          `gorm:"type:int;not null"`
	TopColor  int          `gorm:"type:int;not null"`
	HairColor int          `gorm:"type:int;not null"`
	LegsColor int          `gorm:"type:int;not null"`
	IsNPC     bool         `gorm:"type:boolean;default:0"`
}

// Scan implements the sql.Scanner interface
func (a *Avatar) Scan(value interface{}) error {
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

	avatar := NewAvatarFromString(str)
	if avatar == nil {
		return fmt.Errorf("failed to parse avatar string: %s", str)
	}

	*a = *avatar
	return nil
}

// Value implements the driver.Valuer interface
func (a Avatar) Value() (driver.Value, error) {
	if a.Gender == 0 {
		return nil, nil
	}
	return a.Apperance(), nil
}

func NewAvatarFromString(s string) *Avatar {

	apperances := strings.Split(s, "#")
	var avatar Avatar
	for _, apperance := range apperances {

		// Parse
		values := strings.Split(apperance, "$")
		color, err := strconv.Atoi(values[1])
		if err != nil {
			return nil
		}

		id := values[0][:len(values[0])-1]

		// Set apperance values
		if values[0] == "1001" {
			avatar.Gender = Girl
			avatar.TopColor = color
		} else if values[0] == "1002" {
			avatar.Gender = Boy
			avatar.TopColor = color
		} else if id == "102" {
			avatar.SkinColor = color
		} else if id == "104" {
			avatar.HairColor = color
		} else if id == "105" {
			avatar.LegsColor = color
		}
	}

	return &avatar
}

func (a *Avatar) String() string {
	return fmt.Sprintf("%s+%d+%s", a.Name, a.Gender, a.Apperance())
}

func (a *Avatar) Apperance() string {
	face := fmt.Sprintf("%v$0", 1080+int(a.Gender))
	hat := utils.If(a.IsNPC, "1061$0", "1062$0")
	top := fmt.Sprintf("%v$%v", 1000+int(a.Gender), a.TopColor)
	skin := fmt.Sprintf("%v$%v", 1020+int(a.Gender), a.SkinColor)
	hair := fmt.Sprintf("%v$%v", 1040+int(a.Gender), a.HairColor)
	legs := fmt.Sprintf("%v$%v", 1050+int(a.Gender), a.LegsColor)
	return strings.Join([]string{top, skin, hair, legs, hat, face}, "#")
}

func NewRandomAvatar() Avatar {
	return Avatar{
		Gender:    AvatarGender(rand.Intn(2) + 1),
		SkinColor: rand.Intn(8),
		TopColor:  rand.Intn(17),
		HairColor: rand.Intn(15),
		LegsColor: rand.Intn(17),
		IsNPC:     true,
	}
}
