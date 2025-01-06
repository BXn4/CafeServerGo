package managers

import (
  "fmt"
  "sync"
  "cafego/internal/client"
)


type ClientManager struct {
	mu    sync.Mutex
	clients []*client.Client
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make([]*client.Client, 0),
	}
}

// AddClient adds a new client to the list
func (cm *ClientManager) Add(c *client.Client) {
	cm.mu.Lock() 
	defer cm.mu.Unlock()
	cm.clients = append(cm.clients, c)
}

// RemoveClient removes a client by id 
func (cm *ClientManager) Remove(id int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for i, c := range cm.clients {
		if c.Player.ID == id {
			// Remove client by re-slicing
			cm.clients = append(cm.clients[:i], cm.clients[i+1:]...)
			return
		}
	}
}


// GetClientByID retrieves a client by its ID
func (cm *ClientManager) Get(id int) (*client.Client, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	for _, c := range cm.clients {
		if c.Player.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Cafe with ID %d not found", id)
}
