package managers

import (
	"cafego/internal/client"
	"cafego/internal/database"
	"cafego/internal/objects"
	"sync"
)

type GameManager struct {
	db *database.CafeDB

	locationMutex sync.Mutex
	locations     map[int]*LoadedLocation

	clientMutex sync.Mutex
	clients     []*client.Client
}

func NewGameManager() (*GameManager, error) {

	// Marketplace object
	cafeObj, err := objects.NewMarketplace()
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

	return gm, nil
}
