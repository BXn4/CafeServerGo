package interfaces

import (
	"net"
  "cafego/internal/objects"
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

  // The wrapped cafe object 
  Cafe() *objects.Cafe
}
