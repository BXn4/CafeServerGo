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
	Boy  AvatarGender = 2
)

type AvararWOD struct {
	ID     string `xml:"id,attr"`
	Gender string `xml:"gender,attr"`
	Colors string `xml:"colors,attr"`
}

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

		part := values[0][:len(values[0])]

		if strings.HasPrefix(part, "100") {
			if part == "1001" {
				println(part)
				avatar.Gender = Girl
			}
			if part == "1002" {
				println(part)
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

/*
<wod id="1001" n="Top" g="Avatar" t="Girl" gender="1" colors="798F8F#C5CBE3#76AFB7#61B866#FFE266#F6BCC4#FF7C91#CE8AC5#4A877D#568C49#6BA831#87C7A8#79AEB8#3C7385#B3B4B7#BB8731#BA5C33" />
<wod id="1002" n="Top" g="Avatar" t="Boy" gender="2" colors="475454#C5CBE3#76AFB7#61B866#FFE266#F6BCC4#FF7C91#CE8AC5#4A877D#568C49#6BA831#87C7A8#79AEB8#3C7385#B3B4B7#BB8731#BA5C33" />
<wod id="1021" n="Skin" g="Avatar" t="Girl" gender="1" colors="F4C0A6#F2B187#E89C73#AB6845#CE875C#F3BD90#EEA675#DC9467" />
<wod id="1022" n="Skin" g="Avatar" t="Boy" gender="2" colors="F4C0A6#F2B187#E89C73#AB6845#CE875C#F3BD90#EEA675#DC9467" />
<wod id="1041" n="Hair" g="Avatar" t="Girl" gender="1" colors="FFDC88#FFCC52#D9733B#D68B4A#AF5F3A#73462F#41291C#26120D#DDB25B#BA884C#C77833#9B6134#6C4424#523520#40301C" />
<wod id="1042" n="Hair" g="Avatar" t="Boy" gender="2" colors="FFDC88#FFCC52#D9733B#D68B4A#AF5F3A#73462F#41291C#26120D#DDB25B#BA884C#C77833#9B6134#6C4424#523520#40301C" />
<wod id="1051" n="Legs" g="Avatar" t="Girl" gender="1" colors="475454#C5CBE3#76AFB7#61B866#FFE266#F6BCC4#FF7C91#CE8AC5#4A877D#568C49#6BA831#87C7A8#79AEB8#3C7385#B3B4B7#BB8731#BA5C33" />
<wod id="1052" n="Legs" g="Avatar" t="Boy" gender="2" colors="475454#C5CBE3#76AFB7#61B866#FFE266#F6BCC4#FF7C91#CE8AC5#4A877D#568C49#6BA831#87C7A8#79AEB8#3C7385#B3B4B7#BB8731#BA5C33" />
<wod id="1061" n="Hat" g="Avatar" t="Normal" gender="0" colors="0" />
<wod id="1062" n="Hat" g="Avatar" t="Cook" gender="-1" colors="EEEEEE" />
<wod id="1063" n="Hat" g="Avatar" t="Cookxmas" gender="-1" colors="EEEEEE" />
<wod id="1081" n="Face" g="Avatar" t="Girl" gender="1" colors="0" />
<wod id="1082" n="Face" g="Avatar" t="Boy" gender="2" colors="0" />

func isValidAvatar(s String) bool {

} */
