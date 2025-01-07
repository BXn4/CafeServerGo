package interfaces

import (
	"net"
)

type CafeLocation interface {
	// Add the player to the location by id
	Join(id int, conn net.Conn)

	// Disconnects the player by id
	Leave(id int)

	// Returns the id of the location
	ID() int

	//
	AsResponse() []string

	// Send message to everyone in the location
	Broadcast(arg ...string)

	// Returns the info of the location
	Info() string

	// Returns the fridge
	Fridge() map[int]int

	GetFridgeCapacity() int

	// Returns the owner id
	Owner() int

	//TODO: Edit stuff
}
