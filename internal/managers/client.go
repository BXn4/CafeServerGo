package managers

import (
	"cafego/internal/client"
	"cafego/internal/interfaces"
	"fmt"
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

	newClients := make([]*client.Client, 0, len(gm.clients))
	for i, c := range gm.clients {
		if c.Player == nil && i != id {
			continue
		} else {
			c.SetClientID(i)
			newClients = append(newClients, c)
		}

		for len(c.RequestQueue) > 0 {
			c.RequestQueue <- nil
		}

		if c.Location != nil {
			c.Location.Leave(id)
		}
	}

	gm.clients = newClients

	/* for client := range gm.clients {
	println(client)
	} */
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

func (gm *GameManager) NextClientID() int {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	maxID := 0
	for _, c := range gm.clients {
		if c != nil && c.ClientID > maxID {
			maxID = c.ClientID
		}
	}
	return maxID + 1
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
