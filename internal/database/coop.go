package database

import (
	"cafego/internal/models/coop"
	"time"
)

// CreateCoop creates a new cooperative game session
func (db *CafeDB) CreateCoop(kind int, playerID int, end time.Time) error {
	coop := coop.Coop{
		Host:    playerID,
		Members: []int{playerID}, // Initially only contains the host
		Kind:    kind,
		Dishes:  make(map[int]int),
		Start:   time.Now(),
		End:     end,
	}

	result := db.conn.Create(&coop)
	return result.Error
}

// GetCoop retrieves a cooperative game session by ID
func (db *CafeDB) GetCoop(id int) (coop.Coop, error) {
	var c coop.Coop
	result := db.conn.Where("id = ?", id).First(&c)
	if result.Error != nil {
		return coop.Coop{}, result.Error
	}
	return c, nil
}

// DeleteCoop removes a cooperative game session by ID
func (db *CafeDB) DeleteCoop(id int) error {
	result := db.conn.Delete(&coop.Coop{}, id)
	return result.Error
}

// AddMemberToCoop adds a player to an existing coop session
func (db *CafeDB) AddMemberToCoop(coopID int, playerID int) error {
	coop, err := db.GetCoop(coopID)
	if err != nil {
		return err
	}

	coop.Members = append(coop.Members, playerID)
	result := db.conn.Save(&coop)
	return result.Error
}

// UpdateCoopDishes updates the dishes count in a coop session
func (db *CafeDB) UpdateCoopDishes(coopID int, dishID int, count int) error {
	coop, err := db.GetCoop(coopID)
	if err != nil {
		return err
	}

	coop.Dishes[dishID] = count
	result := db.conn.Save(&coop)
	return result.Error
}
