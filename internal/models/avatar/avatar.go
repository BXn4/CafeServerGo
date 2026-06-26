/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package avatar

import (
	"cafego/internal/utils"
	"database/sql/driver"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type AvatarGender int

const (
	Girl AvatarGender = 1
	Boy  AvatarGender = 2
)

type Avatar struct {
	ID        int          `gorm:"type:int;not null"`
	Name      string       `gorm:"type:string;not null"`
	Gender    AvatarGender `gorm:"type:int;not null"`
	SkinColor int          `gorm:"type:int;not null"`
	TopColor  int          `gorm:"type:int;not null"`
	HairColor int          `gorm:"type:int;not null"`
	LegsColor int          `gorm:"type:int;not null"`
	hat       int          `gorm:"type:int;default:1061"`
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

func (a *Avatar) IsValid() bool {
	if a == nil {
		return false
	}

	if a.SkinColor > validPartColors(1020+int(a.Gender)) || a.SkinColor < 0 ||
		a.TopColor > validPartColors(1000+int(a.Gender)) || a.TopColor < 0 ||
		a.HairColor > validPartColors(1040+int(a.Gender)) || a.HairColor < 0 ||
		a.LegsColor > validPartColors(1050+int(a.Gender)) || a.LegsColor < 0 {
		return false
	}

	return true
}

func NewAvatarFromString(s string) *Avatar {
	// in the request were receiving the new avatar name, and the genders,
	// were just need colors. this causes 0 for color what commes after the gender.
	// [Guest_22776508+2+1042 14] <--- 1042 color value is 0
	parts := strings.Split(s, "+")
	if len(parts) > 0 {
		s = parts[len(parts)-1]
	}

	apperances := strings.Split(s, "#")
	var avatar Avatar
	for _, apperance := range apperances {
		// Parse
		values := strings.Split(apperance, "$")
		color, err := strconv.Atoi(values[1])
		if err != nil {
			return nil
		}

		part := values[0][:len(values[0])]

		if strings.HasPrefix(part, "100") {
			if part == "1001" {
				avatar.Gender = Girl
			}
			if part == "1002" {
				avatar.Gender = Boy
			}
			avatar.TopColor = color
		} else if strings.HasPrefix(part, "102") {
			avatar.SkinColor = color
		} else if strings.HasPrefix(part, "104") {
			avatar.HairColor = color
		} else if strings.HasPrefix(part, "105") {
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
	hat := fmt.Sprintf("%v$0", a.GetAvatarHat())
	top := fmt.Sprintf("%v$%v", 1000+int(a.Gender), a.TopColor)
	skin := fmt.Sprintf("%v$%v", 1020+int(a.Gender), a.SkinColor)
	hair := fmt.Sprintf("%v$%v", 1040+int(a.Gender), a.HairColor)
	legs := fmt.Sprintf("%v$%v", 1050+int(a.Gender), a.LegsColor)
	return strings.Join([]string{top, skin, hair, legs, hat, face}, "#")
}

func (a *Avatar) GetAvatarHat() int {
	if a.hat == 0 {
		if a.IsNPC {
			a.SetAvatarHat(1061)
		} else {
			if utils.GetEventType(time.Now().UTC()) == 3 {
				a.SetAvatarHat(1063)
			} else {
				a.SetAvatarHat(1062)
			}
		}
	}
	return a.hat
}

func (a *Avatar) SetAvatarHat(id int) {
	a.hat = id
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

func validPartColors(id int) int {
	part, err := utils.GetAvatar(id)
	if err != nil {
		return 0
	}
	return len(strings.Split(part.Colors, "#"))
}
