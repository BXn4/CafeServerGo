package managers

import (
	"cafego/internal/database"
	"cafego/internal/objects"
	"fmt"
	"log"
)

func (gm *GameManager) SetCafeDB(db *database.CafeDB) {
	gm.db = db
}

func (gm *GameManager) SetLocation(id int, cafe *LoadedLocation) {
	gm.locationMutex.Lock()
	defer gm.locationMutex.Unlock()

	gm.locations[id] = cafe
}

func (gm *GameManager) RemoveLocation(id int) {
	gm.locationMutex.Lock()
	defer gm.locationMutex.Unlock()

	for i, lc := range gm.locations {
		if lc.Cafe().ID == id {
			// This removes the location by id by not changing the others memory addresses
			gm.db.SaveCafe(lc.cafe)
			fmt.Printf("Saved %v cafe to db\n", lc.cafe.ID)
			delete(gm.locations, i)
			return
		}
	}
}

// NOTE: change this in the future because
// the client and the server wont be able to handle a lot of people in one place
// so we need to cap it
func (gm *GameManager) AddLocation(id int) *LoadedLocation {
	gm.locationMutex.Lock()
	defer gm.locationMutex.Unlock()

	// Get loaded cafe
	item, err := gm.getLocationByID(id)
	if err == nil {
		// If there is a location already loaded return it
		return item
	}

	// If there is no loaded cafe load one
	var cafeObj *objects.Cafe
	cafeObj, err = gm.db.GetCafeByPlayerID(id)
	if err != nil {
		// BIG FUCK UP
		log.Printf("[ERROR] Player with id %v has no cafe in database", id)
		return nil
	}
	cafe := NewLoadedLocation(cafeObj, gm)
	fmt.Printf("Loaded %v cafe from db\n", cafe.cafe.ID)
	gm.locations[id] = cafe

	return cafe
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

func (gm *GameManager) getLocationByID(id int) (*LoadedLocation, error) {
	cafe, ok := gm.locations[id]
	if ok {
		return cafe, nil
	}
	return nil, fmt.Errorf("Cafe with ID %d not found", id)
}
