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

	// This reserves this object so it cannot be interacted with
	// the reservation stays until the reserver unlocks it (like mutex without wait)
	// this should prevent us from iterating over every object in the cafe
	// this returns false is already reserved
	ReserveObject(*objects.CafeObject) bool

	// This returns a reserved object by pos
	GetReservedObject(int, int) *objects.CafeObject

	// This unreserves a dirty table and a chair
	// returns the chair
	GetDirtySpace() *objects.CafeObject

	// This unreserves the reserved object
	UnreserveObject(*objects.CafeObject)

	//
	ClearReservedObjects()

	//
	Owner() (*objects.Player, error)

	//
	GetUniqueCustomerID() int

	//
	AddCustomer(*objects.Customer)

	//
	RemoveCustomer(id int)

	//
	SetRunning(bool)

	//
	IsRunning() bool
}
