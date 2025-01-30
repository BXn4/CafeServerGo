package database

import (
	"cafego/internal/objects"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type CafeDAO struct {
	ID                 int    `gorm:"column:id"`
	PlayerID           int    `gorm:"column:player_id;uniqueIndex"`
	OwnerName          string `gorm:"column:owner_name"`
	Rating             int    `gorm:"column:rating"`
	Luxury             int    `gorm:"column:luxury"`
	Size               int    `gorm:"column:size"`
	Tiles              string `gorm:"column:tiles"`
	Objects            string `gorm:"column:objects"`
	FridgeInventory    string `gorm:"column:fridge_inv"`
	FurnitureInventory string `gorm:"column:furniture_inv"`
	Waiters            string `gorm:"column:waiters"`
}

func (cafeDAO CafeDAO) TableName() string {
	return "cafe"
}

func ConvertCafeDAOToCafe(cafeDAO CafeDAO) (*objects.Cafe, error) {

	// Fill simple cafe object
	cafe := objects.NewCafe(
		cafeDAO.ID,
		cafeDAO.PlayerID,
		cafeDAO.OwnerName,
		cafeDAO.Luxury,
		cafeDAO.Size,
	)
	cafe.SetRating(cafe.GetMinimumRating(cafeDAO.Rating))

	// Parse tiles
	var err error
	err = cafe.ParseTiles(cafeDAO.Tiles)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse tiles: %v", err)
	}

	// Parse objects
	err = cafe.ParseObjectsFromJSON(cafeDAO.Objects)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse objects from json: %v", err)
	}

	// Parse fridge inventory
	cafe.SetFridgeInventory(map[int]int{})
	fridgeInv := strings.Split(cafeDAO.FridgeInventory, "#")
	for _, item := range fridgeInv {

		if item == "" {
			continue
		}

		// Parse item and count
		data := strings.Split(item, "+")
		id, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, fmt.Errorf("Cannot convert data: %v", err)
		}
		count, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot convert count: %v", err)
		}

		// Add to fridge
		cafe.AddToFridge(id, count)
	}

	// Parse furniture inventory
	cafe.SetFurnitureInventory(map[int]int{})
	furnitureInv := strings.Split(cafeDAO.FurnitureInventory, "#")
	for _, item := range furnitureInv {

		// Parse item and count
		data := strings.Split(item, "+")
		id, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, fmt.Errorf("Cannot convert furniture id: %v", err)
		}
		count, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, fmt.Errorf("Cannot convert furniture count: %v", err)
		}

		// Add to furnitures
		cafe.AddFurnitures(id, count)
	}

	// Parse waiters
	waitersRaw := strings.Split(cafeDAO.Waiters, "%")
	for i, waiterRaw := range waitersRaw {
		waiter, err := NewWaiterFromString(waiterRaw)
		if err != nil {
			return nil, err
		}
		waiter.ID = i + 1
		cafe.AddWaiter(waiter)
	}

	return cafe, nil
}

func (db *CafeDB) GetCafeByPlayerID(playerID int) (*objects.Cafe, error) {

	var cafeDAO CafeDAO
	if err := db.conn.Where("player_id = ?", playerID).First(&cafeDAO).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Cafe for player ID %d not found", playerID)
		}
		return nil, fmt.Errorf("Database error: %v", err)
	}
	cafe, err := ConvertCafeDAOToCafe(cafeDAO)
	if err != nil {
		return nil, fmt.Errorf("Cannot convert CafeDAO to Cafe: %v", err)
	}

	return cafe, nil
}

func (db *CafeDB) SaveCafe(cafe *objects.Cafe) error {

	// Build tiles
	var tiles []string
	for i, row := range cafe.GetTiles() {
		for j := range len(row) {
			tiles = append(tiles, strconv.Itoa(cafe.GetTiles()[i][j]))
		}
	}

	// Build objs
	objs := []string{}
	for _, obj := range cafe.GetObjects() {
		objs = append(objs, obj.JSON())
	}

	// Build fridge inventory
	fridgeInv := []string{}
	for i, v := range cafe.GetFridgeInventory() {
		fridgeInv = append(fridgeInv, fmt.Sprintf("%v+%v", i, v))
	}

	// Build furniture inventory
	furnitureInv := []string{}
	for i, v := range cafe.GetFurnitureInventory() {
		furnitureInv = append(furnitureInv, fmt.Sprintf("%v+%v", i, v))
	}

	// Build waiters
	waiters := []string{}
	for _, w := range cafe.GetWaiters() {
		waiters = append(waiters, w.String())
	}

	cafeDAO := CafeDAO{
		ID:                 cafe.GetID(),
		PlayerID:           cafe.GetPlayerID(),
		OwnerName:          cafe.GetOwnerName(),
		Rating:             cafe.GetRating(),
		Luxury:             cafe.GetLuxury(),
		Size:               cafe.GetSize(),
		Tiles:              strings.Join(tiles, "+"),
		Objects:            "[" + strings.Join(objs, ",") + "]",
		FridgeInventory:    strings.Join(fridgeInv, "#"),
		FurnitureInventory: strings.Join(furnitureInv, "#"),
		Waiters:            strings.Join(waiters, "%"),
	}

	if err := db.conn.Save(&cafeDAO).Error; err != nil {
		return fmt.Errorf("Cannot save cafe: %v", err)
	}

	return nil
}
