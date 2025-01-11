package managers

import (
	"cafego/internal/agents"
	"cafego/internal/objects"
	"cafego/internal/types/responses"
	_ "cafego/internal/utils"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// --- LoadedLocation ----------------------------------------------------------
type LoadedLocation struct {
	cafe      *objects.Cafe
	occupants map[int]net.Conn
	mu        sync.Mutex
	gm        *GameManager
	running   bool
}

func NewLoadedLocation(cafe *objects.Cafe, gm *GameManager) *LoadedLocation {
	return &LoadedLocation{
		cafe:      cafe,
		gm:        gm,
		occupants: make(map[int]net.Conn),
		running:   false,
	}
}

func (lc *LoadedLocation) Cafe() *objects.Cafe {
	return lc.cafe
}

// This will send a message to everyone in the loaded cafe
func (lc *LoadedLocation) Broadcast(args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.broadcast(args...)
}

// Same as the Broadcast, just sending to the other players, and not sending it to the source
func (lc *LoadedLocation) Announce(playerID int, args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.announce(playerID, args...)
}

func (lc *LoadedLocation) Join(playerID int, conn net.Conn) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if !lc.running {
		lc.running = true
		go agents.AgentCycle(lc, lc.running)
	}

	// Get joined client
	c, err := lc.gm.GetClient(playerID)
	if err != nil {
		fmt.Printf("Cant get client with id: %v\n", playerID)
		return
	}

	// Set position of player
	c.Player.Position[0] = lc.cafe.PlayerStart[0]
	c.Player.Position[1] = lc.cafe.PlayerStart[1]

	// Send everyone the joined player
	lc.announce(playerID, "juj", "-1", "0", c.Player.String())

	// Add to players in cafe
	println("Added player to occupants: ", playerID)
	lc.occupants[playerID] = conn

	// Send cafe data
	args := []string{"sgc", "-1", "0"}
	args = append(args, lc.cafe.AsResponse()...)
	lc.send(playerID, args...)

	var playersStr []string

	// Get all players in location
	for id, _ := range lc.occupants {
		println("PLAYERS IN CAFE: ", id)
		c, err := lc.gm.GetClient(id)
		if err != nil {
			fmt.Printf("Cant get client with id: %v\n", id)
			return
		}

		playersStr = append(playersStr, c.Player.String())
	}

	args = []string{
		"jul", "-1", "0",
		strings.Join(playersStr, "$"),
	}

	// Send it to joined player
	lc.send(playerID, args...)
}

func (lc *LoadedLocation) Leave(playerID int) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	idStr := strconv.Itoa(playerID)

	lc.announce(playerID, "juq", "-1", "0", idStr)

	delete(lc.occupants, playerID)

	// If there are no players at the location and the location is not the marketplace
	if len(lc.occupants) == 0 && !(lc.cafe.ID < 0) {
		lc.running = false
		// If owner is not online
		if player, _ := lc.gm.GetClient(lc.cafe.PlayerID); player == nil {
			lc.gm.RemoveLocation(lc.cafe.ID) // Delete cafe from manager
		}
	}

}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|
func (lc *LoadedLocation) send(id int, args ...string) {
	msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[SENT TO %v] %s\n", id, msg)
	lc.occupants[id].Write([]byte(msg))
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// This will send a message to everyone in the loaded cafe
func (lc *LoadedLocation) broadcast(args ...string) {

	msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[BROADCAST] %s\n", msg)

	for _, o := range lc.occupants {
		o.Write([]byte(msg))
	}
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// Same as the Broadcast, just sending to the other players, and not sending it to the source
func (lc *LoadedLocation) announce(playerID int, args ...string) {

	msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[ANNOUNCE] %s\n", msg)
	for oid, o := range lc.occupants {
		if oid == playerID {
			continue
		}
		o.Write([]byte(msg))
	}
}
