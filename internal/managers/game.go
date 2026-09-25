/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package managers

import (
	"cafego/internal/client"
	"cafego/internal/database"
	"cafego/internal/models/cafe"
	"cafego/internal/models/event"
	"cafego/internal/models/shop"
	"slices"
	"sync"
	"time"
)

type GameManager struct {
	db *database.CafeDB

	locationMutex sync.Mutex
	locations     map[int]*LoadedLocation

	clientMutex sync.Mutex
	clients     map[int]*client.Client
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
		clients:   make(map[int]*client.Client, 0),
	}

	// Create marketplace
	marketplace := NewLoadedLocation(cafeObj, gm)

	// Add marketplace to cafe list
	gm.SetLocation(-1, marketplace)

	go event.CheckForEvent(10 * time.Minute)
	go shop.CheckForShopAvailablity(10 * time.Minute)

	return gm, nil
}

func (gm *GameManager) SaveAll() error {

	var nonCompletedTutorialPlayers []string // dont need save cafe

	for _, client := range gm.clients {
		if client.Player != nil {
			if !client.Player.GetIsTutorialCompleted() {
				nonCompletedTutorialPlayers = append(nonCompletedTutorialPlayers, client.Player.GetUsername())
				continue
			}
			err := gm.db.SavePlayer(client.Player)
			if err != nil {
				return err
			}
		}
	}

	for _, location := range gm.locations {
		if location != nil {
			if location.cafe.GetRoomType() == cafe.CafeRoom {
				if !slices.Contains(nonCompletedTutorialPlayers, location.Cafe().GetOwnerName()) {
					err := gm.db.SaveCafe(location.cafe)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (gm *GameManager) GetClients() map[int]*client.Client {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	return gm.clients
}
