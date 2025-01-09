package managers

import (
	"cafego/internal/objects"
	_"cafego/internal/utils"
	"cafego/internal/types/responses"
	"net"
	"strconv"
	"sync"
  "fmt"
  "strings"
)

// --- LoadedLocation ----------------------------------------------------------
type LoadedLocation struct {
	cafe           *objects.Cafe
	occupants      map[int]net.Conn
	mu             sync.Mutex
  gm             *GameManager
}



func NewLoadedLocation(cafe *objects.Cafe, gm *GameManager) *LoadedLocation {
	return &LoadedLocation{
		cafe:  cafe,
		gm:  gm,
		occupants:  make(map[int]net.Conn),
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


func (lc *LoadedLocation) Join(playerID int, conn net.Conn) {
	lc.mu.Lock()
  defer lc.mu.Unlock()

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
  lc.broadcast("juj", "-1", "0", c.Player.String())

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

func (lc *LoadedLocation) Leave(id int) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	idStr := strconv.Itoa(id)

	lc.broadcast("juq", "-1", "0", idStr)

	delete(lc.occupants, id)


  // If owner is not online
  if player, _ := lc.gm.GetClient(lc.cafe.PlayerID); player == nil {
    // If there are no players at the location and the location is not the marketplace
    if len(lc.occupants) == 0 && !(lc.cafe.ID < 0) {
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

