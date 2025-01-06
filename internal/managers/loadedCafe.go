package managers

import (
	"cafego/internal/objects"
	"cafego/internal/types/responses"
	"net"
	"strconv"
	"sync"
)

// --- LoadedCafe ----------------------------------------------------------
type LoadedCafe struct {
	cafe      *objects.Cafe
	occupants map[int]net.Conn
	mu        sync.Mutex
}

func NewLoadedCafe(cafe *objects.Cafe) *LoadedCafe {
	return &LoadedCafe{
		cafe:      cafe,
		occupants: make(map[int]net.Conn),
	}
}

func (lc *LoadedCafe) ID() int {
	return lc.cafe.ID
}

func (lc *LoadedCafe) AsResponse() []string {
	return lc.cafe.AsResponse()
}

func (lc *LoadedCafe) Owner() int {
	return lc.cafe.PlayerID
}

func (lc *LoadedCafe) Fridge() map[int]int {
	return lc.cafe.FridgeInventory
}

func (lc *LoadedCafe) Furnitures() map[int]int {
	return lc.cafe.FurnitureInventory
}

func (lc *LoadedCafe) Waiters() []*objects.Waiter {
	return lc.cafe.Waiters
}

func (lc *LoadedCafe) Customers() []*objects.Customer {
	return lc.cafe.Customers
}

// This will send a message to everyone in the loaded cafe
func (lc *LoadedCafe) Broadcast(args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	for _, o := range lc.occupants {
		o.Write([]byte(responses.WrapExtensionResponse(args...)))
	}
}

func (lc *LoadedCafe) Info() string {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	return ""
}

func (lc *LoadedCafe) Join(id int, conn net.Conn) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.occupants[id] = conn

	//args := []string{ "sgc", "-1", "0"}
	//args = append(args, lc.cafe.AsResponse()...)
	//conn.Write([]byte(responses.WrapExtensionResponse(args...)))
}

func (lc *LoadedCafe) Leave(id int) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	idStr := strconv.Itoa(id)

	lc.broadcast("juq", "-1", "0", idStr)

	delete(lc.occupants, id)

	//TODO: Run reverse function provided by cafe manager
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// This will send a message to everyone in the loaded cafe
func (lc *LoadedCafe) broadcast(args ...string) {

	for _, o := range lc.occupants {
		o.Write([]byte(responses.WrapExtensionResponse(args...)))
	}
}
