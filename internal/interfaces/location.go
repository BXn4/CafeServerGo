package interfaces

import (
	"cafego/internal/models/cafe"
	"cafego/internal/models/object"
	"cafego/internal/models/player"
	"cafego/internal/models/simple"
	"cafego/internal/types/responses"
	"time"
)

// This is a wrapper for a cafe
// so we can handle the players inside more easily
type CafeLocation interface {
	// Add the player to the location by id
	Join(id int, channel chan<- responses.Response)

	// Disconnects the player by id
	Leave(id int)

	// Send message to everyone in the location
	Broadcast(arg ...string)

	// Send message to other clients in the location (Not going to send to the source)
	Announce(id int, arg ...string)

	Send(id int, arg ...string)

	// The wrapped cafe object
	Cafe() *cafe.Cafe

	// This reserves this object so it cannot be interacted with
	// the reservation stays until the reserver unlocks it (like mutex without wait)
	// this should prevent us from iterating over every object in the cafe
	// this returns false is already reserved
	ReserveObject(*object.Object) bool

	// This returns a reserved object by pos
	GetReservedObject(simple.Position) *object.Object

	// This unreserves a dirty table and a chair
	// returns the chair
	GetDirtySpace() (*object.Object, *object.Object)

	// This unreserves the reserved object
	UnreserveObject(*object.Object)

	//
	ClearReservedObjects()

	//
	Owner() (*player.Player, error)

	//
	IsEmpty() bool

	AtLocation(int) bool

	//
	GetUniqueCustomerID() int

	//
	SetRunning(bool)

	//
	IsRunning() bool

	//
	GetIsRunning() *bool

	TryStepSleep(time.Duration) bool
}
