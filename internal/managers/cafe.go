package managers

import (
	"cafego/internal/database"
	"cafego/internal/objects"
	"fmt"
	"sync"
)

type CafeManager struct {
	mu    sync.Mutex
	cafes []*LoadedCafe
	db    *database.CafeDB
}

func NewCafeManager() *CafeManager {
	return &CafeManager{
		cafes: make([]*LoadedCafe, 0),
	}
}

func (cm *CafeManager) SetCafeDB(db *database.CafeDB) {
	cm.db = db
}

// RemoveCafe removes a cafe by id
func (cm *CafeManager) Remove(id int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for i, lc := range cm.cafes {
		if lc.ID() == id {
			// This removes the cafe by id by not changing the others memory address
			cm.cafes = append(cm.cafes[:i], cm.cafes[i+1:]...)
			return
		}
	}
}

func (cm *CafeManager) Add(id int) *LoadedCafe {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Get loaded cafe
	item, err := cm.Get(id)
	if err == nil {
		// If there is a cafe already loaded return it
		return item
	}

	// If there is no loaded cafe load one
	var cafeObj *objects.Cafe
	cafeObj, err = cm.db.GetCafeByPlayerID(id)
	if err != nil {
		// BIG FUCK UP
		fmt.Printf("[ERROR] Player with id %v has no cafe in database", id)
		return nil
	}
	cafe := NewLoadedCafe(cafeObj)

	cm.cafes = append(cm.cafes, cafe)

	return cafe
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

func (cm *CafeManager) Get(id int) (*LoadedCafe, error) {

	for _, cafe := range cm.cafes {
		if cafe.ID() == id {
			return cafe, nil
		}
	}
	return nil, fmt.Errorf("Cafe with ID %d not found", id)
}
