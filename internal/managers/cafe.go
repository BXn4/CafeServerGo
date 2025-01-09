package managers

import (
	"cafego/internal/database"
	"cafego/internal/objects"
	"fmt"
	"sync"
)

type CafeManager struct {
	mu    sync.Mutex
	cafes map[int]*LoadedCafe
	db    *database.CafeDB
  clientManager *ClientManager
  marketplace   *LoadedCafe
}

func NewCafeManager(cm *ClientManager) (*CafeManager, error) {
  
  //Marketplace object
  cafeObj, err := objects.NewMarketplace()
  if err != nil {
    return nil, err
  }

  // Create marketplace
  marketplace := NewLoadedCafe(cafeObj, cm, nil)

  // Create loaded cafes map
  cafes := make(map[int]*LoadedCafe, 0)

  // Add marketplace to cafe list
  cafes[-1] = marketplace

	return &CafeManager{
		cafes: cafes,
    clientManager: cm,
    
	},nil
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
      delete(cm.cafes,i)
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
	cafe := NewLoadedCafe(cafeObj, cm.clientManager, cm.Remove)

	cm.cafes[cafe.ID()] = cafe

	return cafe
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

func (cm *CafeManager) Get(id int) (*LoadedCafe, error) {
  cafe, ok := cm.cafes[id]
	if ok {
		return cafe, nil
	}
	return nil, fmt.Errorf("Cafe with ID %d not found", id)
}
