package objects

import (
	"cafego/internal/types/cafetypes"
	"cafego/internal/types/daos"
	"cafego/internal/utils"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
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
	PlayerStart        [2]int
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
	InEditorMode       bool
}

func (c *Cafe) AsResponse() []string {

	var tiles []string
	for i, row := range c.Tiles {
		for j := range len(row) {
			tiles = append(tiles, strconv.Itoa(c.Tiles[i][j]))
		}
	}

	var objs []string
	for _, obj := range c.Objects {
		str := obj.String()
		objs = append(objs, str)
	}

	log.Debugf("OBJECTS: %v", strings.Join(objs, "#"))

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
	cafe.Size = cafe.ExpansionID + 8
	cafe.Tiles = make([][]int, cafe.Size)
	for i := range len(cafe.Tiles) {
		cafe.Tiles[i] = make([]int, cafe.Size)
		for j := range len(cafe.Tiles[i]) {
			value, err := strconv.Atoi(raw_tiles[(i*cafe.Size)+j])
			if err != nil {
				return err
			}
			cafe.Tiles[i][j] = value
		}
	}

	return nil
}

func (cafe *Cafe) ParseObjectsFromJSON(rawObjects string) error {
	var daos []*daos.CafeObjectDAO
	if err := json.Unmarshal([]byte(rawObjects), &daos); err != nil {
		return err
	}
	for _, dao := range daos {

		obj := &CafeObject{
			kind:     dao.Kind,
			pos:      dao.Pos,
			rotation: dao.Rotation,

			dishID:     dao.DishID,
			dishStatus: dao.DishStatus,
			dishAmount: dao.DishAmount,

			fancyIng:   dao.FancyIng,
			startedAt:  dao.StartedAt,
			finishesAt: dao.FinishesAt,
		}

		if obj.IsDoor() {
			cafe.PlayerStart = [2]int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
		}
		cafe.Objects = append(cafe.Objects, obj)
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
			cafe.PlayerStart = [2]int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
		}

		cafe.Objects = append(cafe.Objects, obj)
	}
	return nil
}

func (cafe *Cafe) GetObjectByPos(posX int, posY int) *CafeObject {
	for _, obj := range cafe.Objects {
		if obj.GetPos()[0] == posX && obj.GetPos()[1] == posY {
			return obj
		}
	}
	return nil
}

func (cafe *Cafe) AddNewObject(posX int, posY int, objID int, objRotation int) error {
	obj, err := NewCafeObject(posX, posY, objID, objRotation)
	if err != nil {
		return err
	}
	cafe.Objects = append(cafe.Objects, obj)
	return nil
}
func (cafe *Cafe) RemoveObject(posX int, posY int) error {
	for i, object := range cafe.Objects {
		if object.GetPos()[0] == posX && object.GetPos()[1] == posY {
			cafe.Objects = append(cafe.Objects[:i], cafe.Objects[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object not found at position (%d, %d)", posX, posY)
}

func (c *Cafe) GetPlayerStart() []int {
	for _, obj := range c.Objects {
		if obj.IsDoor() {
			PlayerStart := []int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
			return PlayerStart
		}
	}
	return nil
}

// Returns the tables and the chairs around it
// the chairs should face the right direction
func (c *Cafe) GetEatingSpaces() (tablesAndChairs map[*CafeObject][]*CafeObject) {

	// Get all chairs and tables
	tablesAndChairs = make(map[*CafeObject][]*CafeObject, 0)
	var chairs []*CafeObject
	var tables []*CafeObject
	for _, obj := range c.Objects {
		if obj.IsChair() && obj.GetDishID() == -1 {
			chairs = append(chairs, obj)
		}
		if obj.IsTable() {
			tables = append(tables, obj)
		}
	}

	// Loop through each table and get the chairs facing them
	for _, table := range tables {
		var availableChairs []*CafeObject
		for _, chair := range chairs {

			// Check if char beside the table
			diffX := float64(table.GetPos()[0] - chair.GetPos()[0])
			diffY := float64(table.GetPos()[1] - chair.GetPos()[1])
			if math.Abs(diffX)+math.Abs(diffY) > 1 {
				continue
			}

			// Check if chair faces the table
			facingTable := false
			if chair.GetRotation() == cafetypes.Right && int(diffY) == 1 {
				facingTable = true
			} else if chair.GetRotation() == cafetypes.Left && int(diffY) == -1 {
				facingTable = true
			} else if chair.GetRotation() == cafetypes.Down && int(diffX) == -1 {
				facingTable = true
			} else if chair.GetRotation() == cafetypes.Up && int(diffX) == 1 {
				facingTable = true
			}

			if facingTable {
				availableChairs = append(availableChairs, chair)
			}

		}
		if table == nil || len(availableChairs) == 0 {
			continue
		}
		tablesAndChairs[table] = availableChairs
	}

	return tablesAndChairs
}

func (c *Cafe) UpdateRating(newRating int) {
	rating := c.Rating + newRating

	if rating < 10 {
		rating = 10
	} else if rating > 1000 {
		rating = 1000
	}

	minimumRating := c.GetMinimumRating(rating)

	if rating < minimumRating {
		c.Rating = minimumRating
	} else {
		c.Rating = rating
	}
}

func (c *Cafe) GetMinimumRating(rating int) int {
	minimumRating := min(int((1+0.05*float64(c.Luxury))*10), 500)
	if rating < minimumRating {
		return minimumRating
	}

	return rating
}

func (c *Cafe) GetOldTiles(startX int, startY int, endX int, endY int, tileID int) (int, map[[2]int]int) {
	counts := 0
	oldTiles := make(map[[2]int]int)

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			oldTile := c.Tiles[y][x]

			if oldTile != tileID {
				key := [2]int{y, x}
				oldTiles[key]++
				counts++
			}
		}
	}

	return counts, oldTiles
}
