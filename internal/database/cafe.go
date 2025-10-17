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

func (db *CafeDB) SaveCafe(c *cafe.Cafe) error {
	if c.ID > 0 {
		err := db.conn.Model(&cafe.Cafe{}).
			Where("id = ?", c.ID).
			Updates(map[string]any{
				"rating":        c.GetRating(),
				"luxury":        c.GetLuxury(),
				"expansion_id":  c.GetExpansionID(),
				"tiles":         c.Tiles.String(),
				"objects":       c.Objects.StringForDB(),
				"fridge_inv":    c.FridgeInventory.String(),
				"furniture_inv": c.FurnitureInventory.String(),
				"waiters":       c.Waiters.String(),
			}).Error

		if err != nil {
			return fmt.Errorf("Cant save Cafe: %v", err)
		}
	}

	return nil
}

func (db *CafeDB) UpdateObjects(cafeID int, objects string) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("objects", objects).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateRating(cafeID int, rating int) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("rating", rating).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateLuxury(cafeID int, luxury int) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("luxury", luxury).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateSize(cafeID int, size int) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("size", size).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateExpansionID(cafeID int, expansiondID int) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("expansion_id", expansiondID).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateTiles(cafeID int, tiles string) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("tiles", tiles).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateFridgeInventory(cafeID int, inventory string) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("fridge_inv", inventory).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateFurnitureInventory(cafeID int, inventory string) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("furniture_inv", inventory).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}

func (db *CafeDB) UpdateWaiters(cafeID int, waiters string) error {
	err := db.conn.Model(&cafe.Cafe{}).
		Where("id = ?", cafeID).
		Update("waiters", waiters).Error
	if err != nil {
		return fmt.Errorf("Cant update Cafe: %v", err)
	}

	return nil
}
