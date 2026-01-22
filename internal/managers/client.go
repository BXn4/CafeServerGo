package managers

import (
	"cafego/internal/client"
	"cafego/internal/interfaces"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

// AddClient adds a new client to the list
func (gm *GameManager) AddClient(item interfaces.ManagedItem) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	c := item.(*client.Client)

	clientID := gm.NextClientID()
	c.SetClientID(clientID)

	gm.clients[clientID] = c

	log.Info("Assigned a ID to client: ", clientID)
}

// RemoveClient removes a client by id
func (gm *GameManager) DisconnectClient(id int) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	c, ok := gm.clients[id]
	if !ok {
		log.Info("[Disconnect] Client ID not found:", id)
		return
	}

	log.Info("[Disconnect] Disconnecting client id:", id)

	for len(c.RequestQueue) > 0 {
		select {
		case c.RequestQueue <- nil:
		default:
		}
	}
	time.Sleep(time.Millisecond * 100) // Wait until procceses stop

	for len(c.RequestQueue) > 0 {
		c.RequestQueue <- nil
	}
	time.Sleep(time.Millisecond * 100) // Wait until procceses stop

	delete(gm.clients, id)

	if c.Player != nil {
		if c.Player.GetIsTutorialCompleted() {
			gm.db.SavePlayer(c.Player)
		}

		if c.Location != nil {
			c.Location.Leave(c.Player.GetID())
		}
		c.Player = nil
	}

	log.Info("[Disconnect] Client removed, remaining clients:", len(gm.clients))
}

// GetClient retrieves a client by its ID
func (gm *GameManager) GetClient(id int) (interfaces.ManagedItem, error) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.GetID() == id {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Client with ID %v not found", id)
}

func (gm *GameManager) NextClientID() int {
	for i := 1; ; i++ {
		if _, exists := gm.clients[i]; !exists {
			return i
		}
	}
}

// Checks if client is online
func (gm *GameManager) IsOnline(id int) bool {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	for _, c := range gm.clients {
		if c.Player == nil {
			continue
		}
		if c.Player.GetID() == id {
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
		if c.Player.GetID() == id {
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
		if c.Player.GetID() == id {
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
		if c.Player.GetUsername() == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Client with username %v not found", name)
}
