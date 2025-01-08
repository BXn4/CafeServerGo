package managers

import (
	"cafego/internal/objects"
	"cafego/internal/utils"
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
}

func NewLoadedCafe(cafe *objects.Cafe, cm *ClientManager) *LoadedCafe {
	return &LoadedCafe{
		cafe:  cafe,
		clientManager:  cm,
		occupants:  make(map[int]net.Conn),
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

	for _, o := range lc.occupants {
    msg := responses.WrapExtensionResponse(args...)
	  fmt.Printf("[SENT] %s\n", msg)
		o.Write([]byte(msg))
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

  // Add to players in cafe
	lc.occupants[id] = conn

  // Send cafe data
  args := []string{"sgc", "-1", "0"}
	args = append(args, lc.cafe.AsResponse()...)
  lc.send(id, args...)

  // 

  
  // Send all players in cafe
  for id, _ := range lc.occupants {
    c, err := lc.clientManager.Get(id)
    if err != nil {
      fmt.Printf("Cant get client with id: %v", id)
      return 
    }
    player := c.Player

    params := []string{
      strconv.Itoa(player.ID),
      strconv.Itoa(player.ID),
      strconv.Itoa(player.XP),
      strconv.Itoa(player.Position[0]),
      strconv.Itoa(player.Position[1]),
      strconv.Itoa(player.WorkTimeLeft),
      strconv.Itoa(player.OpenJobs),
      utils.If(player.SeekingJob, "1", "0"),
      utils.If(player.AllowFriendRequests, "1", "0"),
      player.Avatar.String(),
    }

    args := []string{
      "jul", "-1", "0",
      strings.Join(params, "+"),
    }

    lc.send(id, args...)
  }

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
