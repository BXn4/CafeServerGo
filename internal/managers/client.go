package managers

import (
	"cafego/internal/client"
	"cafego/internal/interfaces"
	"fmt"

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

func (gm *GameManager) NextClientID() int {
	for i := 1; ; i++ {
		if _, exists := gm.clients[i]; !exists {
			return i
		}
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
		if c.Player.GetID() == id {
			return c, nil
		}
	}

	return nil, fmt.Errorf("Client with ID %v not found", id)
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

func (gm *GameManager) DisconnectClient(id int) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	c, ok := gm.clients[id]
	if !ok {
		log.Info("[Disconnect] Client ID not found:", id)
		return
	}

	delete(gm.clients, id)

	log.Info("[Disconnect] Disconnecting client id:", id)

	if c.Player != nil {
		if c.Player.GetIsTutorialCompleted() {
			gm.db.SavePlayer(c.Player)
		}

		if c.Location != nil {
			// SAVE -> AFTER LEAVE!!!!!!!!!
			// DONT FLIP IT!!!!!!!!
			if c.Player.GetIsTutorialCompleted() {
				gm.db.SaveCafe(c.Location.Cafe())
			}
			c.Location.Leave(c.Player.GetID())
		}
		c.Player = nil
	}

	log.Info("[Disconnect] Client removed, remaining clients:", len(gm.clients))
}

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
