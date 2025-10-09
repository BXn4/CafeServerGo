package managers

import (
	"cafego/internal/client"
	"cafego/internal/interfaces"
	"fmt"
	"time"
)

// AddClient adds a new client to the list
func (gm *GameManager) AddClient(item interfaces.ManagedItem) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	c := item.(*client.Client)

	gm.clients = append(gm.clients, c)
}

// RemoveClient removes a client by id
func (gm *GameManager) DisconnectClient(id int) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	for i, c := range gm.clients {
		if c.Player == nil {
			continue
		}

		// Send signal to close connection
		for len(c.RequestQueue) > 0 {
			c.RequestQueue <- nil
		}
		time.Sleep(time.Millisecond * 100) // Wait until procceses stop

		// Save player to db
		gm.db.SavePlayer(gm.clients[i].Player)
		c.Player = nil

		// Leave current location
		if c.Location != nil {
			if loc, ok := c.Location.(*LoadedLocation); ok {
				loc.Cancel()
			}
			c.Location.Leave(id)
		}

		// Remove client by re-slicing
		gm.clients = append(gm.clients[:i], gm.clients[i+1:]...)
		return

	}
}

// GetClient retrieves a client by its ID
func (gm *GameManager) GetClient(id int) (interfaces.ManagedItem, error) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.ID == id {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Client with ID %v not found", id)
}

// Checks if client is online
func (gm *GameManager) IsOnline(id int) bool {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.ID == id {
			return true
		}
	}

	return false
}

// Checks if client is online
func (gm *GameManager) isOnline(id int) bool {

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.ID == id {
			return true
		}
	}

	return false
}

// !!! USE IT WITH CARE !!!
// Checks if client is online
func (gm *GameManager) UnsafeIsOnline(id int) bool {

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.ID == id {
			return true
		}
	}

	return false
}

func (gm *GameManager) GetClientByName(name string) (*client.Client, error) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()
	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.Username == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Client with username %v not found", name)
}
