/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package cafe

import (
	"cafego/internal/models/customer"
	"cafego/internal/models/event"
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"cafego/internal/models/waiter"
	"cafego/internal/utils"
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"sync"

	"github.com/charmbracelet/log"
)

const (
	DefaultCafeBackground = 1501
	MarketplaceBackground = 1502
	WinterCafeBackground  = 1503
)

const (
	CafeRoom int = iota
	MarketRoom
)

type Cafe struct {
	ID                 int                        `gorm:"column:id;primaryKey;autoIncrement;type:int"`
	OwnerID            int                        `gorm:"column:owner_id;not null;type:int"`
	OwnerName          string                     `gorm:"column:owner_name;not null;type:text"`
	Rating             int                        `gorm:"column:rating;default:50;type:int"`
	Luxury             int                        `gorm:"column:luxury;default:0;type:int"`
	size               int                        `gorm:"-"`
	background         int                        `gorm:"-"`
	ExpansionID        int                        `gorm:"column:expansion_id;type:int;default:0"`
	Tiles              simple.IntMatrix           `gorm:"column:tiles;type:longtext;not null"`
	Objects            object.ObjectList          `gorm:"column:objects;type:longtext;not null"`
	availableTables    object.ObjectList          `gorm:"-"`
	fridgeCapacity     int                        `gorm:"-"`
	FridgeInventory    simple.IntMap              `gorm:"column:fridge_inv;type:text"`
	FurnitureInventory simple.IntMap              `gorm:"column:furniture_inv;type:text"`
	Waiters            waiter.WaiterList          `gorm:"column:waiters;type:longtext;not null"`
	customers          map[int]*customer.Customer `gorm:"-"`
	playerStart        simple.Position            `gorm:"-"`
	roomType           int                        `gorm:"-"`
	mutex              sync.RWMutex               `gorm:"-"`
}

func (cafe *Cafe) TableName() string {
	return "cafe"
}

func (c *Cafe) AsResponse() []string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(c.OwnerID),
		strconv.Itoa(c.playerStart.X),
		strconv.Itoa(c.playerStart.Y),
		c.OwnerName,
		strconv.Itoa(c.Rating),
		strconv.Itoa(c.Luxury),
		strconv.Itoa(c.ExpansionID),
		strconv.Itoa(c.size),
		strconv.Itoa(c.size),
		strconv.Itoa(c.background),
		c.Tiles.String(),
		c.Objects.String(),
	}
}

func NewCafeForCreation(id, playerID int, name string) *Cafe {
	defaultObjects := *object.ParseObjectList("3+0+901+0#5+0+901+0#5+1+601+3+-1+0#0+2+201+0#5+2+401+0#7+2+601+3+-1+0#7+3+401+0#1+4+351+0#5+4+401+0#1+5+252+0+-1#3+5+301+0+-1+0#5+5+601+1+-1+0#7+5+401+0#1+6+252+0+-1#3+6+301+0+-1+0#7+6+601+1+-1+0#1+7+252+0+-1")

	defaultTiles := simple.IntMatrix{
		{7, 101, 101, 101, 101, 101, 101, 101},
		{101, 4, 4, 4, 4, 4, 4, 4},
		{101, 4, 4, 4, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
	}

	// Players have 1-1 amount already after register: https://youtu.be/8A-BFfhGI5Y?si=8E7NzWJGmJ6_S6NM&t=27
	defaultFridgeInventory := simple.IntMap{1314: 1, 1327: 1}

	defaultStartingWaiter := waiter.GetStartingWaiter()

	return &Cafe{
		ID:              id,
		OwnerID:         playerID,
		OwnerName:       name,
		Objects:         defaultObjects,
		Tiles:           defaultTiles,
		FridgeInventory: defaultFridgeInventory,
		Waiters:         defaultStartingWaiter,
	}
}

func (c *Cafe) Initalize() {
	c.Setsize(c.GetExpansionID())

	totalLuxury := c.CalculateLuxury()
	c.SetLuxury(totalLuxury)

	playerStart := c.FindPlayerStart()
	c.SetPlayerStart(playerStart)

	if event.GetEvent() == 3 {
		c.SetBackground(WinterCafeBackground)
	} else {
		c.SetBackground(DefaultCafeBackground)
	}

	fridgeCapacity := c.CalculateFridgeCapacity()

	c.fridgeCapacity = fridgeCapacity

	c.InitializeCustomers()
}

// ** SETTERS ** // ** SETTERS ** // ** SETTERS ** //
func (c *Cafe) SetRating(newRating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Rating = newRating
}
func (c *Cafe) SetWaiters(waiters []*waiter.Waiter) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Waiters = waiters
}

func (c *Cafe) SetCustomers(customers map[int]*customer.Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.customers = customers
}

func (c *Cafe) AddCustomer(cu *customer.Customer) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.customers[cu.GetID()] = cu
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

func (c *Cafe) RemoveCustomer(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.customers == nil {
		return
	}

	delete(c.customers, id)
}

func (c *Cafe) Setsize(expansionID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	expansionInfo, err := utils.GetExpansion(expansionID)

	if err != nil {
		log.Error("No expansion info found!")
		return
	}

	c.size = expansionInfo.SizeX
}

func (c *Cafe) SetBackground(backgroundID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.background = backgroundID
}

func (c *Cafe) SetLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Luxury = luxury
}

func (c *Cafe) AddLuxury(luxury int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.Luxury += luxury
}

func (c *Cafe) AddRating(rating int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if rating > 0 {
		if c.Rating+rating >= 1000 {
			c.Rating = 1000
			return
		}
	}

	if rating < 0 {
		if c.Rating-rating <= 10 {
			c.Rating = 10
			return
		}
	}

	c.Rating += rating
}

func (c *Cafe) ClearAllCustomers() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.customers = make(map[int]*customer.Customer)
}

func (c *Cafe) CleaAllWaiters() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	for _, w := range c.Waiters {
		w.StopWorking()
	}
	c.Waiters = c.Waiters[:0]
}

func (c *Cafe) InitializeCustomers() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.customers = make(map[int]*customer.Customer)
}

func (c *Cafe) SetTile(x, y, value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Tiles[y][x] = value
}

func (cafe *Cafe) AddNewObject(posX int, posY int, objID int, objRotation int) error {
	cafe.mutex.Lock()
	defer cafe.mutex.Unlock()

	obj, err := object.NewObject(posX, posY, objID, objRotation)
	if err != nil {
		return err
	}

	if obj.IsStove() || obj.IsCounter() {
		obj.SetDishID(-1)
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

func (c *Cafe) SetPlayerStart(pos simple.Position) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.playerStart = pos
}

// ** GETTERS ** // ** GETTERS ** // ** GETTERS ** //
func GetDefaultCafe(id, playerID int, name string) *Cafe {
	defaultObjects := *object.ParseObjectList("3+0+901+0#5+0+901+0#5+1+601+3+-1+0#0+2+201+0#5+2+401+0#7+2+601+3+-1+0#7+3+401+0#1+4+351+0#5+4+401+0#1+5+252+0+-1#3+5+301+0+-1+0#5+5+601+1+-1+0#7+5+401+0#1+6+252+0+-1#3+6+301+0+-1+0#7+6+601+1+-1+0#1+7+252+0+-1")

	defaultTiles := simple.IntMatrix{
		{7, 101, 101, 101, 101, 101, 101, 101},
		{101, 4, 4, 4, 4, 4, 4, 4},
		{101, 4, 4, 4, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
		{101, 1, 1, 1, 4, 4, 4, 4},
	}

	// Players have 1-1 amount already after register: https://youtu.be/8A-BFfhGI5Y?si=8E7NzWJGmJ6_S6NM&t=27
	defaultFridgeInventory := simple.IntMap{1314: 1, 1327: 1}

	defaultStartingWaiter := waiter.GetStartingWaiter()

	return &Cafe{
		ID:              id,
		OwnerID:         playerID,
		OwnerName:       name,
		Objects:         defaultObjects,
		Tiles:           defaultTiles,
		FridgeInventory: defaultFridgeInventory,
		Waiters:         defaultStartingWaiter,
	}
}

func (c *Cafe) GetID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.ID
}

func (c *Cafe) GetOwnerID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.OwnerID
}

func (c *Cafe) GetOwnerName() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.OwnerName
}

func (c *Cafe) GetRating() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	luxury := c.Luxury

	if c.Rating <= 10 {
		c.Rating = 10
	}

	if c.Rating >= 1000 {
		c.Rating = 1000
	}

	minimumRating := min(int((1+0.05*float64(luxury))*10), 1000)

	return max(minimumRating, c.Rating)
}

func (c *Cafe) GetLuxury() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Luxury
}

func (c *Cafe) GetSize() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.size
}

func (c *Cafe) GetBackground() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.background
}

func (c *Cafe) GetExpansionID() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.ExpansionID
}

func (c *Cafe) GetTiles() simple.IntMatrix {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Tiles
}

func (c *Cafe) GetObjects() object.ObjectList {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Objects
}

func (c *Cafe) GetAvailableTables() object.ObjectList {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.availableTables
}

func (c *Cafe) GetWaiters() waiter.WaiterList {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.Waiters
}

func (c *Cafe) GetCustomers() map[int]*customer.Customer {
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

func (c *Cafe) GetPlayerStart() simple.Position {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.roomType == MarketRoom {
		var posX, posY int
		forbidden := map[[2]int]bool{
			{1, 9}:  true, // Object
			{5, 11}: true, // Object
			{6, 5}:  true, // Center object
			{6, 6}:  true, // Center object
			{5, 5}:  true, // Center object
			{5, 6}:  true, // Center object
		}

		for {
			posX = rand.IntN(11) + 1
			posY = rand.IntN(11) + 1

			if !forbidden[[2]int{posX, posY}] {
				break
			}
		}

		return simple.NewPosition(posX, posY)
	}

	return c.playerStart
}

func (c *Cafe) FindPlayerStart() simple.Position {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// the door is in the wall layers, so if its 0, then need to use 1 to get the start position
	doorPosX := utils.If(c.GetDoor().GetPos().X == 0, 1, c.GetDoor().GetPos().X)
	doorPosY := utils.If(c.GetDoor().GetPos().Y == 0, 1, c.GetDoor().GetPos().Y)

	return simple.NewPosition(doorPosX, doorPosY)
}

func (c *Cafe) GetRoomType() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.roomType
}

func (c *Cafe) CalculateLuxury() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var totalLuxury int

	for _, obj := range c.Objects {

		objInfo, err := utils.GetItem(int(obj.GetKind()))
		if err != nil {
			log.Error("Could not find the target object info!")
			continue
		}
		totalLuxury += objInfo.Cash / 4000
		totalLuxury += objInfo.Gold * 2
	}

	return totalLuxury
}

func (c *Cafe) CalculateFridgeCapacity() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var totalCapacity int

	for _, obj := range c.Objects {
		if obj.IsFridge() {
			totalCapacity += 50
		}
	}

	return totalCapacity
}

func (c *Cafe) GetEatingSpaces() (tablesAndChairs map[*object.Object][]*object.Object) {

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	// Get all chairs and tables
	tablesAndChairs = make(map[*object.Object][]*object.Object, 0)
	var chairs []*object.Object
	var tables []*object.Object
	for _, obj := range c.Objects {
		if obj.IsChair() && !obj.GetOccupied() {
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

func (cafe *Cafe) GetDoor() *object.Object {
	cafe.mutex.RLock()
	defer cafe.mutex.RUnlock()

	for _, obj := range cafe.Objects {
		if obj.IsDoor() {
			return obj
		}
	}
	return nil
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
