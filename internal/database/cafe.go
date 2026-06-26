/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package database

import (
	"cafego/internal/models/cafe"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func (db *CafeDB) GetCafeByPlayerID(playerID int) (*cafe.Cafe, error) {

	var c cafe.Cafe
	if err := db.conn.Where("owner_id = ?", playerID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Cafe for player ID %d not found", playerID)
		}
		return nil, fmt.Errorf("Database error: %v", err)
	}

	return &c, nil
}

func (db *CafeDB) SaveCafe(c *cafe.Cafe) error {
	if c.GetRoomType() == cafe.CafeRoom {
		err := db.conn.Model(&cafe.Cafe{}).
			Where("id = ?", c.GetID()).
			Updates(map[string]any{
				"rating":        c.GetRating(),
				"expansion_id":  c.GetExpansionID(),
				"tiles":         c.GetTiles().String(),
				"objects":       c.GetObjects().StringForDB(),
				"fridge_inv":    c.GetFridgeInventory().String(),
				"furniture_inv": c.GetFurnitureInventory().String(),
				"waiters":       c.GetWaiters().String(),
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
