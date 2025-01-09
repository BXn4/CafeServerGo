package objects

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
  "fmt"
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
    fmt.Printf("Error parsing %v to int", data[0])
    return nil, err
  } 

  posY, err := strconv.Atoi(data[1])
  if err != nil {
    fmt.Printf("Error parsing %v to int", data[1])
    return nil, err
  }

  kind, err := strconv.Atoi(data[2])
  if err != nil {
    fmt.Printf("Error parsing %v to int", data[2])
    return nil, err
  }

  rotation, err := strconv.Atoi(data[3])
  if err != nil {
    fmt.Printf("Error parsing %v to int", data[3])
    return nil, err
  }

  cafeObj := CafeObject{
    Pos: []int{posX, posY},
    Kind: CafeObjectKind(kind),
    Rotation: CafeObjectRotation(rotation),
  }

	return &cafeObj, nil
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

func (c *CafeObject) isFridge() bool {
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

			fancyIngStr := "0"
			if c.FancyIng {
				fancyIngStr = "1"
			}
			args = append(args, fancyIngStr)

			if c.StartedAt != nil {
				currentTime := time.Now()
				passedTime := currentTime.Second() - c.StartedAt.Second()
				remainingTime := c.FinishesAt.Second() - currentTime.Second()

				args = append(args, strconv.Itoa(passedTime))
				args = append(args, strconv.Itoa(remainingTime))
			} else {
				args = append(args, "-1")
				args = append(args, "-1")
			}
		}
		// if COUNTER
	} else if c.IsCounter() {
		args = append(args, strconv.Itoa(c.DishID))
		args = append(args, strconv.Itoa(c.DishAmount))
		// if CHAIR
	} else if c.IsChair() {
		args = append(args, strconv.Itoa(c.DishID))
		args = append(args, strconv.Itoa(c.DishStatus))
		// if VENDING
	} else if c.IsVending() {
		args = append(args, strconv.Itoa(c.DishID))
		args = append(args, strconv.Itoa(c.DishAmount))
	}

	return strings.Join(args, "+")
}
