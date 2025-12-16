package managers

import (
	"cafego/internal/database"
	"cafego/internal/models/cafe"
	"cafego/internal/models/customer"
	"cafego/internal/models/event"
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
			owner := lc.Cafe().GetOwnerName()
			player, err := gm.db.GetPlayerByName(owner)
			if err == nil {
				if player.GetXP() > 0 {
					gm.db.SaveCafe(lc.Cafe())
					log.Debugf("Saved %v cafe to db", lc.cafe.GetID())
				}
			}

			if lc.Cafe().AgentCycleBinded {
				lc.Cafe().AgentCycleBinded = false
				for _, c := range lc.Cafe().GetCustomers() {
					delete(lc.Cafe().Customers, c.GetID())
				}

				for _, w := range lc.Cafe().Waiters {
					w.StopWorking()
				}

				lc.Cafe().Waiters = lc.Cafe().Waiters[:0]
			}
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

	if event.GetEvent() == 3 {
		cafeObj.Background = cafe.WinterCafeBackground
	} else {
		cafeObj.Background = cafe.DefaultCafeBackground
	}

	//
	loc := NewLoadedLocation(cafeObj, gm)
	gm.locations[id] = loc

	loc.Cafe().Customers = make(map[int]*customer.Customer)

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
