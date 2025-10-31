package managers

import (
	"cafego/internal/client"
	"cafego/internal/database"
	"cafego/internal/models/cafe"
	"sync"
	"time"
)

type GameManager struct {
	db *database.CafeDB

	locationMutex sync.Mutex
	locations     map[int]*LoadedLocation

	clientMutex sync.Mutex
	clients     []*client.Client

	gameEvent              int
	gameEventDaysLeft      int
	unavailableIngredients []int
}

func NewGameManager() (*GameManager, error) {

	// Marketplace object
	cafeObj, err := cafe.NewMarketplace(-1) // Default marketplace
	if err != nil {
		return nil, err
	}

	// Create game manager
	gm := &GameManager{
		locations: make(map[int]*LoadedLocation, 0),
		clients:   make([]*client.Client, 0),
	}

	// Create marketplace
	marketplace := NewLoadedLocation(cafeObj, gm)

	// Add marketplace to cafe list
	gm.SetLocation(-1, marketplace)

	go gm.CheckForEvent(10 * time.Minute)
	go gm.CheckForShopAvailablity()

	return gm, nil
}

func (gm *GameManager) SaveAll() error {

	for _, client := range gm.clients {
		if client.Player != nil {
			err := gm.db.SavePlayer(client.Player)
			if err != nil {
				return err
			}
		}
	}

	for _, location := range gm.locations {
		if location != nil {
			err := gm.db.SaveCafe(location.cafe)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
