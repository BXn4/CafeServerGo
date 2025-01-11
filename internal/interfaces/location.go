package interfaces

import (
	"cafego/internal/objects"
	"net"
)

// This is a wrapper for a cafe
// so we can handle the players inside more easily
type CafeLocation interface {
	// Add the player to the location by id
	Join(id int, conn net.Conn)

	// Disconnects the player by id
	Leave(id int)

	// Send message to everyone in the location
	Broadcast(arg ...string)

	// Send message to other clients in the location (Not going to send to the source)
	Announce(id int, arg ...string)

	// The wrapped cafe object
	Cafe() *objects.Cafe
}
