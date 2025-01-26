package managers

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/objects"
	"cafego/internal/types/responses"
	_ "cafego/internal/utils"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/log"
)

// --- LoadedLocation ----------------------------------------------------------
type LoadedLocation struct {
	cafe         *objects.Cafe
	occupants    map[int](chan<- responses.Response)
	mu           sync.Mutex
	gm           *GameManager
	running      bool
	reservedObjs []*objects.CafeObject
}

func NewLoadedLocation(cafe *objects.Cafe, gm *GameManager) *LoadedLocation {
	return &LoadedLocation{
		cafe:      cafe,
		gm:        gm,
		occupants: make(map[int](chan<- responses.Response), 0),
		running:   false,
	}
}

// TODO: Change this its shiti design and throws nil pointer exception
// when more clients want to use it
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

func (lc *LoadedLocation) IsEmpty() bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return len(lc.occupants) == 0
}

func (lc *LoadedLocation) IsRunning() bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.running
}

func (lc *LoadedLocation) GetIsRunning() *bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return &lc.running
}

func (lc *LoadedLocation) SetRunning(b bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.running = b

	// Change waiter switch
	for _, w := range lc.cafe.GetWaiters() {
		w.IsWorking = b
	}

	// Start agent cycle
	if b {
		go agents.AgentCycle(lc)
	}
}

func (lc *LoadedLocation) setRunning(b bool) {
	lc.running = b

	// Change waiter switch
	for _, w := range lc.cafe.GetWaiters() {
		w.IsWorking = b
	}

	// Start agent cycle
	if b {
		go agents.AgentCycle(lc)
	}
}

func (lc *LoadedLocation) Join(playerID int, channel chan<- responses.Response) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	log.Printf("--------------------------")

	if !lc.running {
		lc.running = true
		go agents.AgentCycle(lc)
	}

	// Get joined client
	c, err := lc.gm.GetClient(playerID)
	if err != nil {
		log.Printf("Cant get client with id: %v\n", playerID)
		return
	}

	// Set position of player
	c.(*client.Client).Player.Position = lc.cafe.GetPlayerStart()

	// Send everyone the joined player
	lc.announce(playerID, "juj", "-1", "0", c.(*client.Client).Player.String())

	// Add to players in cafe
	lc.occupants[playerID] = channel

	// Send cafe data
	args := []string{"sgc", "-1", "0"}
	args = append(args, lc.cafe.AsResponse()...)
	lc.send(playerID, args...)

	var playersStr []string

	// Get all players in location
	for id := range lc.occupants {
		c, err := lc.gm.GetClient(id)
		if err != nil {
			log.Printf("Cant get client with id: %v\n", id)
			return
		}

		playersStr = append(playersStr, c.(*client.Client).Player.String())
	}

	// Send it to joined player
	lc.send(playerID, "jul", "-1", "0", strings.Join(playersStr, "$"))
}

// Leaves the location and broadcasts leave to every one
func (lc *LoadedLocation) Leave(id int) {

	lc.mu.Lock()
	defer lc.mu.Unlock()

	idStr := strconv.Itoa(id)

	// Send leave to everyone
	lc.announce(id, "juq", "-1", "0", idStr)

	// Delete from occupants
	delete(lc.occupants, id)

	// If there are no players at the location and the location is not the marketplace
	if len(lc.occupants) == 0 && !(lc.cafe.GetID() < 0) {
		lc.setRunning(false)
	}

}

func (lc *LoadedLocation) ClearReservedObjects() {
	for _, obj := range lc.reservedObjs {
		if obj.IsChair() {
			obj.SetDishID(-1)
		}
	}
	lc.reservedObjs = []*objects.CafeObject{}
}

// Returns the first table and unreserves it
func (lc *LoadedLocation) GetDirtySpace() *objects.CafeObject {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	chairIndex := -1
	for i, o := range lc.reservedObjs {
		if o.IsChair() && o.GetDishID() == -2 {
			chairIndex = i
		}
	}
	if chairIndex == -1 {
		return nil
	}
	chair := lc.reservedObjs[chairIndex]

	// Get associated table
	for j, o := range lc.reservedObjs {
		if !o.IsTable() {
			continue
		}
		nr := chair.GetNormalizedRotation()
		if o.GetPos()[0] == chair.GetPos()[0]+nr[0] && o.GetPos()[1] == chair.GetPos()[1]+nr[1] {
			// Unreserve chair
			lc.reservedObjs = append(lc.reservedObjs[:chairIndex], lc.reservedObjs[chairIndex+1:]...)
			if chairIndex < j {
				j--
			}

			// Unreserve table
			lc.reservedObjs = append(lc.reservedObjs[:j], lc.reservedObjs[j+1:]...)

			return chair
		}
	}

	return nil
}

// Get reserved item by pos
func (lc *LoadedLocation) GetReservedObject(x, y int) *objects.CafeObject {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	for _, o := range lc.reservedObjs {
		if o.GetPos()[0] == x && o.GetPos()[1] == y {
			return o
		}
	}

	return nil
}

func (lc *LoadedLocation) ReserveObject(obj *objects.CafeObject) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Check if object is reserved
	for _, o := range lc.reservedObjs {
		if o.GetPos() == obj.GetPos() {
			return false
		}
	}

	// Add to reserved
	lc.reservedObjs = append(lc.reservedObjs, obj)

	return true
}

func (lc *LoadedLocation) UnreserveObject(obj *objects.CafeObject) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Search object index
	index := -1
	for i, o := range lc.reservedObjs {
		if o.GetPos() == obj.GetPos() {
			index = i
			break
		}
	}

	// Check for bug (this should not happen if the code is rigth)
	if index == -1 {
		return
	}

	// Delete reservation
	lc.reservedObjs = append(lc.reservedObjs[:index], lc.reservedObjs[index+1:]...)
}

func (lc *LoadedLocation) GetUniqueCustomerID() int {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	var ids []int
	for _, customer := range lc.cafe.GetCustomers() {
		ids = append(ids, customer.GetID())
	}

	id := 101
	for slices.Contains(ids, id) {
		id++
	}
	return id
}

func (lc *LoadedLocation) AddCustomer(customer *objects.Customer) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.Cafe().AddCustomer(customer)
}

func (lc *LoadedLocation) Owner() (*objects.Player, error) {
	c, err := lc.gm.GetClient(lc.cafe.GetID())
	if err != nil {
		return nil, err
	}
	return c.(*client.Client).Player, nil
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|
func (lc *LoadedLocation) send(id int, args ...string) {
	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-3), "to %v: %s", id, resp.Wrap())
	lc.occupants[id] <- resp
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// This will send a message to everyone in the loaded cafe
func (lc *LoadedLocation) broadcast(args ...string) {

	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-1), "%s", resp.Wrap())
	for _, channel := range lc.occupants {
		channel <- resp
	}
}

// |========================================|
// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// |========================================|

// Same as the Broadcast, just sending to the other players, and not sending it to the source
func (lc *LoadedLocation) announce(playerID int, args ...string) {

	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-2), "%s", resp.Wrap())
	for id, channel := range lc.occupants {
		if id == playerID {
			continue
		}
		channel <- resp
	}
}
