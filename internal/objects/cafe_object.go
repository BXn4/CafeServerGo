package objects

import (
	"cafego/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CafeObjectKind int

// TODO: Write down types
const (
	STOVE   CafeObjectKind = 0
	COUNTER                = 1
	CHAIR                  = 2
	TABLE                  = 3
	VENDING                = 4
	OTHER                  = 5
)

type CafeObjectRotation int

const (
	Up CafeObjectRotation = iota
	Left
	Down
	Right
)

type CafeObject struct {
	Kind     CafeObjectKind     `json:"id"`
	Pos      []int              `json:"pos"`
	Rotation CafeObjectRotation `json:"rotation"`

	DishID     int `json:"dish_id,omitempty"`
	DishStatus int `json:"dish_status,omitempty"`
	DishAmount int `json:"dish_amount,omitempty"`

	FancyIng   bool       `json:"fancy_ing,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishesAt *time.Time `json:"finishes_at,omitempty"`
}

func NewCafeObject(posX int, posY int, objID int, objRotation int) (*CafeObject, error) {
	cafeObj := CafeObject{
		Pos:      []int{posX, posY},
		Kind:     CafeObjectKind(objID),
		Rotation: CafeObjectRotation(objRotation),
	}
	return &cafeObj, nil
}

func NewCafeObjectFromJSON(s string) (*CafeObject, error) {

	var cafeObj CafeObject

	if err := json.Unmarshal([]byte(s), &cafeObj); err != nil {
		return nil, err
	}

	return &cafeObj, nil
}

func NewCafeObjectFromString(s string) (*CafeObject, error) {

	data := strings.Split(s, "+")

	posX, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v to int: %v", data[0], err)
	}

	posY, err := strconv.Atoi(data[1])
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v to int: %v", data[1], err)
	}

	kind, err := strconv.Atoi(data[2])
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v to int: %v", data[2], err)
	}

	rotation, err := strconv.Atoi(data[3])
	if err != nil {
		return nil, fmt.Errorf("Error parsing %v to int: %v", data[3], err)
	}

	cafeObj := CafeObject{
		Pos:      []int{posX, posY},
		Kind:     CafeObjectKind(kind),
		Rotation: CafeObjectRotation(rotation),
	}

	return &cafeObj, nil
}

func (c *CafeObject) IsWall() bool {
	return 101 <= c.Kind && c.Kind <= 135
}

func (c *CafeObject) IsDoor() bool {
	return 201 <= c.Kind && c.Kind <= 207
}

func (c *CafeObject) IsStove() bool {
	return 252 <= c.Kind && c.Kind <= 259
}

func (c *CafeObject) IsCounter() bool {
	return 301 <= c.Kind && c.Kind <= 312
}

func (c *CafeObject) IsChair() bool {
	return 601 <= c.Kind && c.Kind <= 624
}

func (c *CafeObject) IsTable() bool {
	return 401 <= c.Kind && c.Kind <= 423
}

func (c *CafeObject) IsVending() bool {
	return c.Kind == 1701
}

func (c *CafeObject) IsFridge() bool {
	return 351 <= c.Kind && c.Kind <= 358
}

func (c *CafeObject) String() string {
	args := []string{
		strconv.Itoa(c.Pos[0]),
		strconv.Itoa(c.Pos[1]),
		strconv.Itoa(int(c.Kind)),
		strconv.Itoa(int(c.Rotation)),
	}

	// if STOVE
	if c.IsStove() {

		args = append(args, strconv.Itoa(c.DishID))
		if !(c.DishID == -1 || c.DishID == -2) {

			fancyIngStr := utils.If(c.FancyIng, "1", "0")
			args = append(args, fancyIngStr)

			if c.StartedAt != nil {
				currentTime := time.Now().UTC()
				passedTime := currentTime.Sub(*c.StartedAt).Seconds()
				remainingTime := c.FinishesAt.Sub(currentTime).Seconds()

				args = append(args,
					strconv.Itoa(int(passedTime)),
					strconv.Itoa(int(remainingTime)),
				)
			} else {
				args = append(args, "-1", "-1")
			}
		}
		// if COUNTER
	} else if c.IsCounter() {
		args = append(args, strconv.Itoa(c.DishID), strconv.Itoa(c.DishAmount))
		// if CHAIR
	} else if c.IsChair() {
		args = append(args, strconv.Itoa(c.DishID), strconv.Itoa(c.DishStatus))
		// if VENDING
	} else if c.IsVending() {
		args = append(args, strconv.Itoa(c.DishID), strconv.Itoa(c.DishAmount))
	}

	return strings.Join(args, "+")
}

func (c *CafeObject) JSON() string {
	b, err := json.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(b)
}

func (c *CafeObject) GetNormalizedRotation() [2]int {
	if c.Rotation == Up {
		return [2]int{1, 0}
	} else if c.Rotation == Down {
		return [2]int{-1, 0}
	} else if c.Rotation == Left {
		return [2]int{0, -1}
	} else /* Right */ {
		return [2]int{0, 1}
	}
}
