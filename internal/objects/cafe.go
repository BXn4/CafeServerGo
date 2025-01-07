package objects

import (
	"strconv"
	"strings"
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
		string(DefaultBackground),
		strings.Join(tiles, "+"),
		strings.Join(objs, "#"),
	}
	return args
}

func (c *Cafe) GetFridgeCapacity() int {
	fridgeCount := 0

	for _, obj := range c.Objects {
		if obj.isFridge() {
			fridgeCount++
		}
	}
	return fridgeCount * 50
}
