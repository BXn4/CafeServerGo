package database

import (
	"cafego/internal/objects"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type CafeDAO struct {
	ID                 int    `json:"id" form:"id" gorm:"column:id"`
	PlayerID           int    `json:"player_id" form:"player_id" gorm:"column:player_id"`
	OwnerName          string `json:"owner_name" form:"owner_name" gorm:"column:owner_name"`
	Rating             int    `json:"rating" form:"rating" gorm:"column:rating"`
	Luxury             int    `json:"luxury" form:"luxury" gorm:"column:luxury"`
	ExpansionID        int    `json:"expansion_id" form:"expansion_id" gorm:"column:expansion_id"`
	Tiles              string `json:"tiles" form:"tiles" gorm:"column:tiles"`
	Objects            string `json:"objects" form:"objects" gorm:"column:objects"`
	FridgeInventory    string `json:"fridge_inv" form:"fridge_inv" gorm:"column:fridge_inv"`
	FurnitureInventory string `json:"furniture_inv" form:"furniture_inv" gorm:"column:furniture_inv"`
	Waiters            string `json:"waiters" form:"waiters" gorm:"column:waiters"`
}

func ConvertCafeDAOToCafe(cafeDAO CafeDAO) (*objects.Cafe, error) {

	// Fill simple cafe object
	var cafe objects.Cafe
	cafe.ID = cafeDAO.ID
	cafe.PlayerID = cafeDAO.PlayerID
	cafe.Rating = cafeDAO.Rating
	cafe.Luxury = cafeDAO.Luxury
	cafe.ExpansionID = cafeDAO.ExpansionID
	cafe.OwnerName = cafeDAO.OwnerName

	// Parse tiles
	raw_tiles := strings.Split(cafeDAO.Tiles, "+")
	size := cafeDAO.ExpansionID + 8
	cafe.Tiles = make([][]int, size)
	for i, _ := range cafe.Tiles {
		cafe.Tiles[i] = make([]int, size)
		for j, _ := range cafe.Tiles[i] {
			value, err := strconv.Atoi(raw_tiles[(i*size)+j])
			if err != nil {
				return nil, err
			}
			cafe.Tiles[i][j] = value
		}
	}

	// Parse objects
	var objs []objects.CafeObject
	if err := json.Unmarshal([]byte(cafeDAO.Objects), &objs); err != nil {
		return nil, err
	}
	for _, obj := range objs {
		if obj.IsDoor() {
			cafe.PlayerStart = []int{
				obj.Pos[0],
				obj.Pos[1],
			}
			if cafe.PlayerStart[0] == 0 {
				cafe.PlayerStart[0] = 1
			}
			if cafe.PlayerStart[1] == 0 {
				cafe.PlayerStart[1] = 1
			}
		}
		cafe.Objects = append(cafe.Objects, &obj)
	}

	// Parse fridge inventory
	cafe.FridgeInventory = map[int]int{}
	raw_frInv := strings.Split(cafeDAO.FridgeInventory, "#")
	for _, item := range raw_frInv {

		// Parse item and count
		data := strings.Split(item, "+")
		id, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, err
		}

		// Add to fridge
		cafe.FridgeInventory[id] = count
	}

	// Parse furniture inventory
	cafe.FurnitureInventory = map[int]int{}
	raw_fuInv := strings.Split(cafeDAO.FridgeInventory, "#")
	for _, item := range raw_fuInv {
		// Parse item and count
		data := strings.Split(item, "+")
		id, err := strconv.Atoi(data[0])
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, err
		}

		// Add to furnitures
		cafe.FurnitureInventory[id] = count
	}

	// Parse waiters
	var daos []WaiterDAO
	if err := json.Unmarshal([]byte(cafeDAO.Waiters), &daos); err != nil {
		return nil, err
	}
	for _, waiterDAO := range daos {
		waiter, err := ConvertWaiterDAOToWaiter(&waiterDAO)
		if err != nil {
			return nil, err
		}
		cafe.Waiters = append(cafe.Waiters, waiter)
	}

	return &cafe, nil
}

func (db *CafeDB) GetCafeByPlayerID(player_id int) (*objects.Cafe, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.conn.QueryRow("SELECT * FROM cafe WHERE player_id = ?", player_id)

	var cafeDAO CafeDAO
	err := row.Scan(
		&cafeDAO.ID,
		&cafeDAO.PlayerID,
		&cafeDAO.Rating,
		&cafeDAO.Luxury,
		&cafeDAO.ExpansionID,
		&cafeDAO.Tiles,
		&cafeDAO.Objects,
		&cafeDAO.OwnerName,
		&cafeDAO.FridgeInventory,
		&cafeDAO.FurnitureInventory,
		&cafeDAO.Waiters,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("ID NOT FOUND")
		}
		fmt.Errorf("SQL ERR: %v", err)
		return nil, err
	}

	cafe, err := ConvertCafeDAOToCafe(cafeDAO)
	if err != nil {
		return nil, err
	}

	return cafe, nil
}
