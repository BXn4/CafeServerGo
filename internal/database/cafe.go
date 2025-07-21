package database

import (
	"cafego/internal/models/cafe"
	"errors"
	"fmt"

	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func (db *CafeDB) GetCafeByPlayerID(playerID int) (*cafe.Cafe, error) {

	var c cafe.Cafe
	if err := db.conn.Where("player_id = ?", playerID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Cafe for player ID %d not found", playerID)
		}
		return nil, fmt.Errorf("Database error: %v", err)
	}

	c.Background = cafe.DefaultBackground
	c.GetPlayerStart()

	return &c, nil
}

func (db *CafeDB) SaveCafe(cafe *cafe.Cafe) error {

	// If it is a marketplace dont save it
	if cafe.ID < 0 {
		return nil
	}

	// Delete temporary info from objects
	for _, obj := range cafe.Objects {
		if obj.IsChair() {
			obj.SetDishID(0)
			obj.SetDishStatus(0)
		}
	}

	// Save to database
	if err := db.conn.Save(&cafe).Error; err != nil {
		return fmt.Errorf("Cannot save cafe: %v", err)
	}

	return nil
}
