package managers

import (
	"cafego/internal/client"
	"fmt"
)


// AddClient adds a new client to the list
func (gm *GameManager) AddClient(c *client.Client) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()

	gm.clients = append(gm.clients, c)
}

// RemoveClient removes a client by id
func (gm *GameManager) DisconnectClient(id int) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()
	for i, c := range gm.clients {
    if c.Player == nil { continue; }
		if c.Player.ID == id {
			// Remove client by re-slicing
			gm.clients = append(gm.clients[:i], gm.clients[i+1:]...)
      c.Disconnect()
			return
		}
	}
}

// RemoveClient removes a client by name
func (gm *GameManager) DisconnectClientByName(name string) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()
	for i, c := range gm.clients {
    if c.Player == nil { continue; }
		if c.Player.Username == name {
			// Remove client by re-slicing
			gm.clients = append(gm.clients[:i], gm.clients[i+1:]...)
      c.Disconnect()
			return
		}
	}
}

// GetClientByID retrieves a client by its ID
func (gm *GameManager) GetClient(id int) (*client.Client, error) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()
	for _, c := range gm.clients {
    if c.Player == nil { continue; }
		if c.Player.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Cafe with ID %d not found", id)
}

func (gm *GameManager) GetClientByName(name string) (*client.Client, error) {
	gm.clientMutex.Lock()
	defer gm.clientMutex.Unlock()
	for _, c := range gm.clients {
    if c.Player == nil { continue; }
		if c.Player.Username == name {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Cafe with username %v not found", name)
}
