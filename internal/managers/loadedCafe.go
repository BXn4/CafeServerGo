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

// --- LoadedCafe ----------------------------------------------------------
type LoadedCafe struct {
	cafe           *objects.Cafe
	occupants      map[int]net.Conn
	mu             sync.Mutex
  clientManager  *ClientManager
  deleteFunc     func(int)  // This is a loopback function provided by the cafe manager
}

func NewLoadedCafe(cafe *objects.Cafe, cm *ClientManager, deleteFunc func(int)) *LoadedCafe {
	return &LoadedCafe{
		cafe:  cafe,
		clientManager:  cm,
		occupants:  make(map[int]net.Conn),
    deleteFunc: deleteFunc,
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

func (lc *LoadedCafe) GetFridgeMaxCapacity() int {
	return lc.cafe.GetFridgeMaxCapacity()
}

func (lc *LoadedCafe) GetFridgeFreeSpace() int {
	return lc.cafe.GetFridgeFreeSpace()
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

  lc.broadcast(args...)
}

// This will send a message to everyone in the loaded cafe except one with id
func (lc *LoadedCafe) BroadcastExcept(id int, args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

  lc.broadcastExcept(id, args...)
}

func (lc *LoadedCafe) Info() string {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	return ""
}

func (lc *LoadedCafe) Join(id int, conn net.Conn) {
	lc.mu.Lock()
  defer lc.mu.Unlock()

  // Get joined client
  c, err := lc.clientManager.Get(id)
  if err != nil {
      fmt.Printf("Cant get client with id: %v\n", id)
      return 
  }

  // Set position of player
  c.Player.Position[0] = lc.cafe.PlayerStart[0]
  c.Player.Position[1] = lc.cafe.PlayerStart[1]

  // Send everyone the joined player
  lc.broadcast("juj", "-1", "0", c.Player.String())

  // Add to players in cafe
	lc.occupants[id] = conn

  // Send cafe data
  args := []string{"sgc", "-1", "0"}
	args = append(args, lc.cafe.AsResponse()...)
  lc.send(id, args...)

  // Get all players in cafe
  var playersStr []string
  
  // Get all players in cafe
  for id, _ := range lc.occupants {
    c, err := lc.clientManager.Get(id)
    if err != nil {
      fmt.Printf("Cant get client with id: %v\n", id)
      return 
    }
    player := c.Player

    playersStr = append(playersStr, player.String())
  }

  args = []string{
    "jul", "-1", "0",
    strings.Join(playersStr, "$"),
  }
  
  // Send it to joined player
  lc.send(id, args...)
}

func (lc *LoadedCafe) Leave(id int) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	idStr := strconv.Itoa(id)

	lc.broadcast("juq", "-1", "0", idStr)

	delete(lc.occupants, id)


  // If owner is not online
  if player, _ := lc.clientManager.Get(lc.cafe.ID); player == nil {
    // If there are no players at the location
    if len(lc.occupants) == 0 {
      lc.mu.Unlock()
      lc.deleteFunc(lc.cafe.ID) // Delete cafe from manager
      println("TEST DOES THIS WORK?????")
    }
  }

}



// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|
func (lc *LoadedCafe) send(id int, args ...string) {
  msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[SENT TO %v] %s\n", id, msg)
	lc.occupants[id].Write([]byte(msg))
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// This will send a message to everyone in the loaded cafe
func (lc *LoadedCafe) broadcast(args ...string) {
  
  msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[BROADCAST] %s\n", msg)

	for _, o := range lc.occupants {
		o.Write([]byte(msg))
	}
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|
func (lc *LoadedCafe) broadcastExcept(id int, args ...string) {
  
  msg := responses.WrapExtensionResponse(args...)
	fmt.Printf("[BROADCAST] %s\n", msg)

	for oid, o := range lc.occupants {
    if oid == id { continue; }
		o.Write([]byte(msg))
	}
}
