package objects

import (
	"fmt"
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

func (a *Avatar) String() string {
	return fmt.Sprintf("%s+%d+%s", a.Name, a.Gender, a.Apperance())
}

func (a *Avatar) Apperance() string {

	face := fmt.Sprintf("%v$0", 1080+int(a.Gender))

	hat := "1062$0"
	if a.IsNPC {
		hat = "1061$0"
	}

	top := fmt.Sprintf("%v$%v", 1000+int(a.Gender), a.TopColor)
	skin := fmt.Sprintf("%v$%v", 1020+int(a.Gender), a.SkinColor)
	hair := fmt.Sprintf("%v$%v", 1040+int(a.Gender), a.HairColor)
	legs := fmt.Sprintf("%v$%v", 1050+int(a.Gender), a.LegsColor)

	return strings.Join([]string{top, skin, hair, legs, hat, face}, "#")
}
