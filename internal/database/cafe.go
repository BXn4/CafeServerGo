package database

import (
	"cafego/internal/objects"
	"cafego/internal/types/daos"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"

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
	cafe := objects.NewCafe(
		cafeDAO.ID,
		cafeDAO.PlayerID,
		cafeDAO.OwnerName,
		cafeDAO.Luxury,
		cafeDAO.ExpansionID,
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
	raw_frInv := strings.Split(cafeDAO.FridgeInventory, "#")
	for _, item := range raw_frInv {

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
	raw_fuInv := strings.Split(cafeDAO.FurnitureInventory, "#")
	for _, item := range raw_fuInv {

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
	var daos []daos.WaiterDAO
	if err := json.Unmarshal([]byte(cafeDAO.Waiters), &daos); err != nil {
		return nil, err
	}
	for _, waiterDAO := range daos {
		waiter, err := ConvertWaiterDAOToWaiter(&waiterDAO)
		if err != nil {
			return nil, err
		}
		cafe.AddWaiter(waiter)
	}

	return cafe, nil
}

func (db *CafeDB) GetCafeByPlayerID(player_id int) (*objects.Cafe, error) {

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
			return nil, fmt.Errorf("ID NOT FOUND")
		}
		return nil, fmt.Errorf("\n\tSQL ERR: %v", err)
	}

	cafe, err := ConvertCafeDAOToCafe(cafeDAO)
	if err != nil {
		return nil, fmt.Errorf("\n\tCannot convert CafeDAO to Cafe: %v", err)
	}

	return cafe, nil
}

func (db *CafeDB) SaveCafe(cafe *objects.Cafe) {

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
		waiters = append(waiters, w.JSON())
	}

	result, err := db.conn.Exec(
		" UPDATE cafe SET "+
			"rating = ?,"+
			"luxury = ?,"+
			"expansion_id = ?,"+
			"tiles = ?,"+
			"objects = ?,"+
			"fridge_inv = ?,"+
			"furniture_inv = ?,"+
			"waiters = ? "+
			"WHERE player_id = ?",
		cafe.GetRating(),
		cafe.GetLuxury(),
		cafe.GetExpansionID(),
		strings.Join(tiles, "+"),
		"["+strings.Join(objs, ", ")+"]",
		strings.Join(fridgeInv, "#"),
		strings.Join(furnitureInv, "#"),
		"["+strings.Join(waiters, ", ")+"]",
		cafe.GetID(),
	)

	if err != nil {
		log.Errorf("Cant save cafe: %v\n", err)
		return
	}

	// Check how many rows were affected
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("Error fetching rows affected:", err)
	}

}
