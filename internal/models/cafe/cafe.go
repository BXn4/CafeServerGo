package cafe

import (
	"cafego/internal/models/customer"
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"cafego/internal/models/waiter"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
	"sync"
)

type CafeBackground string

const (
	DefaultBackground     CafeBackground = "1501"
	MarketplaceBackground CafeBackground = "1502"
	WinterBackground      CafeBackground = "1503"
)

type Cafe struct {
	ID                 int                  `gorm:"primaryKey;autoIncrement"`
	PlayerID           int                  `gorm:"not null;type:int"`
	OwnerName          string               `gorm:"not null;type:text"`
	Rating             int                  `gorm:"default:50"`
	Luxury             int                  `gorm:"default:0"`
	Size               int                  `gorm:"default:8"`
	Background         CafeBackground       `gorm:"type:text"`
	Tiles              simple.IntMatrix     `gorm:"type:longtext"`
	Objects            object.ObjectList    `gorm:"type:longtext"`
	availableTables    object.ObjectList    `gorm:"-"`
	fridgeCapacity     int                  `gorm:"-"`
	FridgeInventory    simple.IntMap        `gorm:"column:fridge_inv;type:text"`
	FurnitureInventory simple.IntMap        `gorm:"column:furniture_inv;type:text"`
	Waiters            waiter.WaiterList    `gorm:"type:longtext;"`
	customers          []*customer.Customer `gorm:"-"`
	playerStart        *simple.Position     `gorm:"-"`
	mutex              sync.RWMutex         `gorm:"-"`
}

func (cafe *Cafe) TableName() string {
	return "cafe"
}

func NewCafe(id int, playerID int, ownerName string, luxury int, size int) *Cafe {
	return &Cafe{
		ID:         id,
		PlayerID:   playerID,
		OwnerName:  ownerName,
		Luxury:     luxury,
		Size:       size,
		Background: DefaultBackground,
	}
}

func NewCafeForCreation(id, playerID int, name string) *Cafe {
	return &Cafe{
		ID:        id,
		PlayerID:  playerID,
		OwnerName: name,
		Objects:   *object.ParseObjectList("3+0+901+0#5+0+901+0#5+1+601+3+-1+0#0+2+201+0#5+2+401+0#7+2+601+3+-1+0#7+3+401+0#1+4+351+0#5+4+401+0#1+5+252+0+-1#3+5+301+0+-1+0#5+5+601+1+-1+0#7+5+401+0#1+6+252+0+-1#3+6+301+0+-1+0#7+6+601+1+-1+0#1+7+252+0+-1"),
		Tiles: simple.IntMatrix{
			{7, 101, 101, 101, 101, 101, 101, 101, 101},
			{4, 4, 4, 4, 4, 4, 4, 101},
			{4, 4, 4, 4, 4, 4, 4, 101},
			{4, 4, 4, 4, 4, 4, 4, 101},
			{1, 1, 1, 4, 4, 4, 4, 101},
			{1, 1, 1, 4, 4, 4, 4, 101},
			{1, 1, 1, 4, 4, 4, 4, 101},
			{1, 1, 1, 4, 4, 4, 4},
		}}
}

func (c *Cafe) AsResponse() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.playerStart == nil {
		println("NO PLAYERSTART FOUND!") // NILOOLASODAOSDOL was here!
		c.getPlayerStart()
		println(c.playerStart == nil)
	}

	return []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(c.PlayerID),
		strconv.Itoa(c.playerStart.X),
		strconv.Itoa(c.playerStart.Y),
		c.OwnerName,
		strconv.Itoa(c.Rating),
		strconv.Itoa(c.Luxury),
		strconv.Itoa(c.Size - 8),
		strconv.Itoa(c.Tiles.Size()),
		strconv.Itoa(c.Tiles.Size()),
		string(c.Background),
		c.Tiles.String(),
		c.Objects.String(),
	}
}

func (cafe *Cafe) GetObjectByPos(pos simple.Position) *object.Object {
	cafe.mutex.RLock()
	defer cafe.mutex.RUnlock()
	for _, obj := range cafe.Objects {
		if obj.GetPos() == pos {
			return obj
		}
	}
	return nil
}

func (cafe *Cafe) GetObjectByPosXY(x, y int) *object.Object {
	cafe.mutex.RLock()
	defer cafe.mutex.RUnlock()
	pos := simple.NewPosition(x, y)
	for _, obj := range cafe.Objects {
		if obj.GetPos() == pos {
			return obj
		}
	}
	return nil
}

func (cafe *Cafe) AddNewObject(posX int, posY int, objID int, objRotation int) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

	obj, err := object.NewObject(posX, posY, objID, objRotation)
	if err != nil {
		return err
	}
	cafe.Objects = append(cafe.Objects, obj)
	return nil
}

func (c *Cafe) RemoveObject(pos simple.Position) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for i, obj := range c.Objects {
		if obj.GetPos() == pos {
			c.Objects = append(c.Objects[:i], c.Objects[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("object not found at position %v", pos)
}

func (c *Cafe) GetPlayerStart() simple.Position {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.getPlayerStart()
}

func (c *Cafe) getPlayerStart() simple.Position {

	// If we already have one dont search
	if c.playerStart != nil {
		return *c.playerStart
	}

	// If market
	if c.ID < 0 {
		pos := simple.NewPosition(1, 2)
		c.playerStart = &pos
		return pos
	}

	for _, obj := range c.Objects {
		if obj.IsDoor() {
			PlayerStart := simple.NewPosition(
				utils.If(obj.GetPos().X == 0, 1, obj.GetPos().X),
				utils.If(obj.GetPos().Y == 0, 1, obj.GetPos().Y),
			)
			c.playerStart = &PlayerStart
			return PlayerStart
		}
	}

	return simple.NewPosition(1, 2)
}

// Returns the tables and the chairs around it
// the chairs should face the right direction
func (c *Cafe) GetEatingSpaces() (tablesAndChairs map[*object.Object][]*object.Object) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Get all chairs and tables
	tablesAndChairs = make(map[*object.Object][]*object.Object, 0)
	var chairs []*object.Object
	var tables []*object.Object
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
		var availableChairs []*object.Object
		for _, chair := range chairs {

			// Check if char beside the table
			diffX := float64(table.GetPos().X - chair.GetPos().X)
			diffY := float64(table.GetPos().Y - chair.GetPos().Y)
			if math.Abs(diffX)+math.Abs(diffY) > 1 {
				continue
			}

			// Check if chair faces the table
			facingTable := false
			if chair.GetRotation() == object.Right && int(diffY) == 1 {
				facingTable = true
			} else if chair.GetRotation() == object.Left && int(diffY) == -1 {
				facingTable = true
			} else if chair.GetRotation() == object.Down && int(diffX) == -1 {
				facingTable = true
			} else if chair.GetRotation() == object.Up && int(diffX) == 1 {
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
	rating := c.Rating + newRating

	if rating < 10 {
		rating = 10
	} else if rating > 1000 {
		rating = 1000
	}

	minimumRating := c.GetMinimumRating(rating)
	c.Rating = utils.If(rating < minimumRating, minimumRating, rating)
}

func (c *Cafe) GetMinimumRating(rating int) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	minimumRating := min(int((1+0.05*float64(c.Luxury))*10), 500)
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

// Getters
func (c *Cafe) GetID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.ID
}

func (c *Cafe) GetPlayerID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.PlayerID
}

func (c *Cafe) GetOwnerName() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.OwnerName
}

func (c *Cafe) GetRating() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Rating
}

func (c *Cafe) GetLuxury() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Luxury
}

func (c *Cafe) GetSize() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Size
}

func (c *Cafe) GetBackground() CafeBackground {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Background
}

func (c *Cafe) GetTiles() [][]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Tiles
}

func (c *Cafe) GetObjects() []*object.Object {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Objects
}

func (c *Cafe) GetFridgeInventory() simple.IntMap {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.FridgeInventory
}

func (c *Cafe) GetFurnitureInventory() map[int]int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.FurnitureInventory
}

func (c *Cafe) GetWaiters() []*waiter.Waiter {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.Waiters
}

func (c *Cafe) GetCustomers() []*customer.Customer {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.customers
}

func (c *Cafe) GetCustomer(id int) *customer.Customer {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for _, cs := range c.customers {
		if cs.GetID() == id {
			return cs
		}
	}

	return nil
}

// Setters
func (c *Cafe) SetID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.ID = id
}

func (c *Cafe) SetPlayerID(playerID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.PlayerID = playerID
}

func (c *Cafe) SetOwnerName(ownerName string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.OwnerName = ownerName
}

func (c *Cafe) SetRating(rating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Rating = rating
}

func (c *Cafe) AddRating(rating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Rating += rating
}

func (c *Cafe) SetLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Luxury = luxury
}

func (c *Cafe) AddLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Luxury = luxury
}

func (c *Cafe) SetSize(size int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Size = size
}

func (c *Cafe) SetBackground(background CafeBackground) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Background = background
}

func (c *Cafe) SetTiles(tiles [][]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Tiles = tiles
}

func (c *Cafe) SetObjects(objects []*object.Object) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Objects = objects
}

func (c *Cafe) SetFridgeCapacity(fridgeCapacity int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.fridgeCapacity = fridgeCapacity
}

func (c *Cafe) SetFridgeInventory(fridgeInventory map[int]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.FridgeInventory = fridgeInventory
}

func (c *Cafe) SetFurnitureInventory(furnitureInventory map[int]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.FurnitureInventory = furnitureInventory
}

func (c *Cafe) SetWaiters(waiters []*waiter.Waiter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Waiters = waiters
}

func (c *Cafe) SetCustomers(customers []*customer.Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.customers = customers
}

func (c *Cafe) AddCustomer(customer *customer.Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.customers = append(c.customers, customer)
}

func (c *Cafe) AddWaiter(waiter *waiter.Waiter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Waiters = append(c.Waiters, waiter)
}

func (c *Cafe) RemoveWaiter(waiterID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	index := -1
	for i, waiter := range c.Waiters {
		if waiter.GetID() == waiterID {
			waiter.StopWorking()
			index = i
			break
		}
	}
	if index == -1 {
		return
	}
	c.Waiters = append(c.Waiters[:index], c.Waiters[index+1:]...)
}

func (c *Cafe) SetTile(x, y, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Tiles[y][x] = value
}

// Removes a customer from the list by id
func (c *Cafe) RemoveCustomer(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Search index
	index := -1
	for i, customer := range c.customers {
		if customer.GetID() == id {
			index = i
		}
	}

	// If not found return
	if index == -1 {
		return
	}

	// Remove from slice
	c.customers = append(c.customers[:index], c.customers[index+1:]...)
}
