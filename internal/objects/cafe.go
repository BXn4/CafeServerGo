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
	"sync"
)

type CafeBackground string

const (
	DefaultBackground     CafeBackground = "1501"
	MarketplaceBackground                = "1502"
	WinterBackground                     = "1503"
)

type Cafe struct {
	id                 int
	playerID           int
	playerStart        [2]int
	ownerName          string
	rating             int
	luxury             int
	expansionID        int
	size               int
	background         CafeBackground
	tiles              [][]int
	objects            []*CafeObject
	availableTables    []*CafeObject
	fridgeCapacity     int
	fridgeInventory    map[int]int
	furnitureInventory map[int]int
	waiters            []*Waiter
	customers          []*Customer
	inEditorMode       bool
	mutex              sync.RWMutex
}

func NewCafe(id int, playerID int, ownerName string, luxury int, expansionID int) *Cafe {
	return &Cafe{
		id:          id,
		playerID:    playerID,
		ownerName:   ownerName,
		luxury:      luxury,
		expansionID: expansionID,
		background:  DefaultBackground,
	}
}

func (c *Cafe) AsResponse() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	var tiles []string
	for i, row := range c.tiles {
		for j := range len(row) {
			tiles = append(tiles, strconv.Itoa(c.tiles[i][j]))
		}
	}

	var objs []string
	for _, obj := range c.objects {
		str := obj.String()
		objs = append(objs, str)
	}

	args := []string{
		strconv.Itoa(c.id),
		strconv.Itoa(c.playerID),
		strconv.Itoa(c.playerStart[0]),
		strconv.Itoa(c.playerStart[1]),
		c.ownerName,
		strconv.Itoa(c.rating),
		strconv.Itoa(c.luxury),
		strconv.Itoa(c.expansionID),
		strconv.Itoa(len(c.tiles)),
		strconv.Itoa(len(c.tiles[0])),
		string(c.background),
		strings.Join(tiles, "+"),
		strings.Join(objs, "#"),
	}
	return args
}

func (cafe *Cafe) ParseTiles(rawTiles string) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

	// Parse tiles
	raw_tiles := strings.Split(rawTiles, "+")
	cafe.size = cafe.expansionID + 8
	cafe.tiles = make([][]int, cafe.size)
	for i := range len(cafe.tiles) {
		cafe.tiles[i] = make([]int, cafe.size)
		for j := range len(cafe.tiles[i]) {
			value, err := strconv.Atoi(raw_tiles[(i*cafe.size)+j])
			if err != nil {
				return err
			}
			cafe.tiles[i][j] = value
		}
	}

	return nil
}

func (cafe *Cafe) ParseObjectsFromJSON(rawObjects string) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

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
			cafe.playerStart = [2]int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
		}
		cafe.objects = append(cafe.objects, obj)
	}
	return nil
}

func (cafe *Cafe) ParseObjects(rawObjects string) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

	cafe.objects = []*CafeObject{}

	objsStr := strings.Split(rawObjects, "#")

	for _, objStr := range objsStr {

		obj, err := NewCafeObjectFromString(objStr)

		if err != nil {
			return err
		}

		if obj.IsDoor() {
			cafe.playerStart = [2]int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
		}

		cafe.objects = append(cafe.objects, obj)
	}
	return nil
}

func (cafe *Cafe) GetObjectByPos(posX int, posY int) *CafeObject {
	cafe.mutex.RLock()
	defer cafe.mutex.RUnlock()
	for _, obj := range cafe.objects {
		if obj.GetPos()[0] == posX && obj.GetPos()[1] == posY {
			return obj
		}
	}
	return nil
}

func (cafe *Cafe) AddNewObject(posX int, posY int, objID int, objRotation int) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

	obj, err := NewCafeObject(posX, posY, objID, objRotation)
	if err != nil {
		return err
	}
	cafe.objects = append(cafe.objects, obj)
	return nil
}
func (cafe *Cafe) RemoveObject(posX int, posY int) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()
	for i, object := range cafe.objects {
		if object.GetPos()[0] == posX && object.GetPos()[1] == posY {
			cafe.objects = append(cafe.objects[:i], cafe.objects[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object not found at position (%d, %d)", posX, posY)
}

func (c *Cafe) GetPlayerStart() [2]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	for _, obj := range c.objects {
		if obj.IsDoor() {
			PlayerStart := [2]int{
				utils.If(obj.GetPos()[0] == 0, 1, obj.GetPos()[0]),
				utils.If(obj.GetPos()[1] == 0, 1, obj.GetPos()[1]),
			}
			return PlayerStart
		}
	}
	return [2]int{1, 1}
}

// Returns the tables and the chairs around it
// the chairs should face the right direction
func (c *Cafe) GetEatingSpaces() (tablesAndChairs map[*CafeObject][]*CafeObject) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Get all chairs and tables
	tablesAndChairs = make(map[*CafeObject][]*CafeObject, 0)
	var chairs []*CafeObject
	var tables []*CafeObject
	for _, obj := range c.objects {
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
	c.mutex.Lock()
	defer c.mutex.Unlock()
	rating := c.rating + newRating

	if rating < 10 {
		rating = 10
	} else if rating > 1000 {
		rating = 1000
	}

	minimumRating := c.GetMinimumRating(rating)
	c.rating = utils.If(rating < minimumRating, minimumRating, rating)
}

func (c *Cafe) GetMinimumRating(rating int) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	minimumRating := min(int((1+0.05*float64(c.luxury))*10), 500)
	if rating < minimumRating {
		return minimumRating
	}

	return rating
}

func (c *Cafe) GetOldTiles(startX int, startY int, endX int, endY int, tileID int) (int, map[[2]int]int) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	counts := 0
	oldTiles := make(map[[2]int]int)

	for y := startY; y <= endY; y++ {
		for x := startX; x <= endX; x++ {
			oldTile := c.tiles[y][x]

			if oldTile != tileID {
				key := [2]int{y, x}
				oldTiles[key]++
				counts++
			}
		}
	}

	return counts, oldTiles
}

// Getters
func (c *Cafe) GetID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.id
}

func (c *Cafe) GetPlayerID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.playerID
}

func (c *Cafe) GetOwnerName() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.ownerName
}

func (c *Cafe) GetRating() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.rating
}

func (c *Cafe) GetLuxury() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.luxury
}

func (c *Cafe) GetExpansionID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.expansionID
}

func (c *Cafe) GetSize() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.size
}

func (c *Cafe) GetBackground() CafeBackground {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.background
}

func (c *Cafe) GetTiles() [][]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.tiles
}

func (c *Cafe) GetObjects() []*CafeObject {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.objects
}

func (c *Cafe) GetFridgeInventory() map[int]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.fridgeInventory
}

func (c *Cafe) GetFurnitureInventory() map[int]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.furnitureInventory
}

func (c *Cafe) GetWaiters() []*Waiter {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.waiters
}

func (c *Cafe) GetCustomers() []*Customer {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.customers
}

func (c *Cafe) InEditorMode() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.inEditorMode
}

// Setters
func (c *Cafe) SetID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id = id
}

func (c *Cafe) SetPlayerID(playerID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.playerID = playerID
}

func (c *Cafe) SetOwnerName(ownerName string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ownerName = ownerName
}

func (c *Cafe) SetRating(rating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.rating = rating
}

func (c *Cafe) AddRating(rating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.rating += rating
}

func (c *Cafe) SetLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.luxury = luxury
}

func (c *Cafe) AddLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.luxury = luxury
}

func (c *Cafe) SetExpansionID(expansionID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.expansionID = expansionID
}

func (c *Cafe) SetSize(size int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.size = size
}

func (c *Cafe) SetBackground(background CafeBackground) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.background = background
}

func (c *Cafe) SetTiles(tiles [][]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.tiles = tiles
}

func (c *Cafe) SetObjects(objects []*CafeObject) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.objects = objects
}

func (c *Cafe) SetFridgeCapacity(fridgeCapacity int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.fridgeCapacity = fridgeCapacity
}

func (c *Cafe) SetFridgeInventory(fridgeInventory map[int]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.fridgeInventory = fridgeInventory
}

func (c *Cafe) SetFurnitureInventory(furnitureInventory map[int]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.furnitureInventory = furnitureInventory
}

func (c *Cafe) SetWaiters(waiters []*Waiter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waiters = waiters
}

func (c *Cafe) SetCustomers(customers []*Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.customers = customers
}

func (c *Cafe) AddCustomer(customer *Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.customers = append(c.customers, customer)
}

func (c *Cafe) SetInEditorMode(inEditorMode bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.inEditorMode = inEditorMode
}

func (c *Cafe) AddWaiter(waiter *Waiter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.waiters = append(c.waiters, waiter)
}

func (c *Cafe) RemoveWaiter(waiterID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	index := -1
	for i, waiter := range c.waiters {
		if waiter.ID == waiterID {
			waiter.StopWorking()
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.waiters = append(c.waiters[:index], c.waiters[index+1:]...)
}

func (c *Cafe) SetTile(x, y, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.tiles[y][x] = value
}

func (c *Cafe) RemoveCustomer(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	index := -1
	for i, customer := range c.customers {
		if customer.GetID() == id {
			index = i
		}
	}

	if index == -1 {
		return
	}
	c.customers = append(c.customers[:index], c.customers[index+1:]...)
}

func (c *Cafe) SetPlayerStart(playerStart [2]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.playerStart = playerStart
}
