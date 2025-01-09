package objects

import (
	"strconv"
	"strings"
	"encoding/json"
  "cafego/internal/utils"
)

type CafeBackground string

const (
	DefaultBackground     CafeBackground = "1501"
	MarketplaceBackground                = "1502"
	WinterBackground                     = "1503"
)

type Cafe struct {
	ID                 int
	PlayerID           int
	PlayerStart        []int
	OwnerName          string
	Rating             int
	Luxury             int
	ExpansionID        int
	Size               int
  Background         CafeBackground
	Tiles              [][]int
	Objects            []*CafeObject
	AvailableTables    []*CafeObject
	FridgeCapacity     int
	FridgeInventory    map[int]int
	FurnitureInventory map[int]int
	Waiters            []*Waiter
	Customers          []*Customer
	CustomerCycle      string //TODO: Make it 'Task'
}

func (c *Cafe) AsResponse() []string {

	var tiles []string
	for i, row := range c.Tiles {
		for j, _ := range row {
			tiles = append(tiles, strconv.Itoa(c.Tiles[i][j]))
		}
	}

	var objs []string
	for _, v := range c.Objects {
		objs = append(objs, v.String())
	}

	args := []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(c.PlayerID),
		strconv.Itoa(c.PlayerStart[0]),
		strconv.Itoa(c.PlayerStart[1]),
		c.OwnerName,
		strconv.Itoa(c.Rating),
		strconv.Itoa(c.Luxury),
		strconv.Itoa(c.ExpansionID),
		strconv.Itoa(len(c.Tiles)),
		strconv.Itoa(len(c.Tiles[0])),
		string(c.Background),
		strings.Join(tiles, "+"),
		strings.Join(objs, "#"),
	}
	return args
}

func (cafe *Cafe) ParseTiles(rawTiles string) error {
  // Parse tiles
	raw_tiles := strings.Split(rawTiles, "+")
  size := cafe.ExpansionID + 8
  cafe.Tiles = make([][]int, size)
	for i, _ := range cafe.Tiles {
		cafe.Tiles[i] = make([]int, size)
		for j, _ := range cafe.Tiles[i] {
			value, err := strconv.Atoi(raw_tiles[(i*size)+j])
			if err != nil {
				return err
			}
			cafe.Tiles[i][j] = value
		}
	}

  return nil
}

func (cafe *Cafe) ParseObjectsFromJSON(rawObjects string) error {
	var objs []CafeObject
	if err := json.Unmarshal([]byte(rawObjects), &objs); err != nil {
		return err
	}
	for _, obj := range objs {
		if obj.IsDoor() {
			cafe.PlayerStart = []int{
				utils.If(obj.Pos[0] == 0, 1, obj.Pos[0]),
				utils.If(obj.Pos[1] == 0, 1, obj.Pos[1]),
			}
		}
		cafe.Objects = append(cafe.Objects, &obj)
	}
  return nil
}

func (cafe *Cafe) ParseObjects(rawObjects string) error {
  cafe.Objects = []*CafeObject{}

  objsStr := strings.Split(rawObjects, "#")

	for _, objStr := range objsStr {

    obj, err := NewCafeObjectFromString(objStr)

    if err != nil {
      return err
    }

		if obj.IsDoor() {
			cafe.PlayerStart = []int{
				utils.If(obj.Pos[0] == 0, 1, obj.Pos[0]),
				utils.If(obj.Pos[1] == 0, 1, obj.Pos[1]),
			}
		}

		cafe.Objects = append(cafe.Objects, obj)
	}
  return nil
}


func (c *Cafe) GetFridgeMaxCapacity() int {
	fridgeCount := 0

	for _, obj := range c.Objects {
		if obj.isFridge() {
			fridgeCount++
		}
	}
	return fridgeCount * 50
}

func (c *Cafe) GetFridgeFreeSpace() int {
	freeSpace := c.GetFridgeMaxCapacity()

	for ingredientID := range c.FridgeInventory {
		// Fancy does not take space in the fridge inventory. Fancys starts at ID 1401
		if ingredientID < 1401 {
			freeSpace -= 1
		}
	}

	return freeSpace
}
