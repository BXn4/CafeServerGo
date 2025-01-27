package objects

import (
	"cafego/internal/utils"
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
	Name      string
	Gender    AvatarGender
	SkinColor int
	TopColor  int
	HairColor int
	LegsColor int
	IsNPC     bool
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

		// Set values
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
		Name:      "Customer",
		Gender:    AvatarGender(rand.Intn(2) + 1),
		SkinColor: rand.Intn(8),
		TopColor:  rand.Intn(17),
		HairColor: rand.Intn(15),
		LegsColor: rand.Intn(17),
		IsNPC:     true,
	}
}
