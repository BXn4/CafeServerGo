package managers

import (
	"cafego/internal/database"
	"cafego/internal/models/cafe"
	"fmt"

	"github.com/charmbracelet/log"
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

	println("REMOVED CAFE: ", id)

	for i, lc := range gm.locations {
		if lc.cafe.GetID() == id {
			// This removes the location by id by not changing the others memory addresses
			gm.db.SaveCafe(lc.Cafe())
			log.Debugf("Saved %v cafe to db", lc.cafe.GetID())
			delete(gm.locations, i)
			return
		}
	}
}

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
	var cafeObj *cafe.Cafe
	cafeObj, err = gm.db.GetCafeByPlayerID(id)
	if err != nil {
		log.Errorf("Player with id %v has no cafe in database: %v", id, err)
		return nil
	}

	if gm.GetEvent() == 3 {
		cafeObj.Background = cafe.WinterBackground
	} else {
		cafeObj.Background = cafe.DefaultBackground
	}

	//
	loc := NewLoadedLocation(cafeObj, gm)
	gm.locations[id] = loc

	println("LOADED CAFE: ", id)

	return loc
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
