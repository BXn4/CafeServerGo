package managers

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/models/cafe"
	"cafego/internal/models/customer"
	"cafego/internal/models/object"
	"cafego/internal/models/player"
	"cafego/internal/models/simple"
	"cafego/internal/models/waiter"
	"cafego/internal/types/responses"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

// --- LoadedLocation ----------------------------------------------------------
type LoadedLocation struct {
	cafe         *cafe.Cafe
	occupants    map[int](chan<- responses.Response)
	mu           sync.Mutex
	gm           *GameManager
	running      bool
	reservedObjs []*object.Object
}

func NewLoadedLocation(cafe *cafe.Cafe, gm *GameManager) *LoadedLocation {
	return &LoadedLocation{
		cafe:      cafe,
		gm:        gm,
		occupants: make(map[int](chan<- responses.Response), 0),
		running:   false,
	}

}

// TODO: Change this its shit design and throws nil pointer exception
// when more clients want to use it
// INSTEAD we should create a threadsafe interface for cafe
func (lc *LoadedLocation) Cafe() *cafe.Cafe {
	return lc.cafe
}

// This will send a message to everyone in the loaded cafe
func (lc *LoadedLocation) Broadcast(args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.broadcast(args...)
}

// This will send a message to a user in the location
func (lc *LoadedLocation) Send(id int, args ...string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.send(id, args...)
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

func (lc *LoadedLocation) AtLocation(id int) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	_, ok := lc.occupants[id]
	return ok
}

func (lc *LoadedLocation) IsRunning() bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.running
}

// TODO: Change this
func (lc *LoadedLocation) GetIsRunning() *bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return &lc.running
}

func (lc *LoadedLocation) SetRunning(b bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.setRunning(b)
}

func (lc *LoadedLocation) setRunning(b bool) {

	lc.running = b
}

func (lc *LoadedLocation) Join(playerID int, channel chan<- responses.Response) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

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

	// Send cafe data
	args := []string{"sgc", "-1", "0"}
	args = append(args, lc.cafe.AsResponse()...)
	lc.send(playerID, args...)

	// Send it to joined player
	lc.send(playerID, "jul", "-1", "0", strings.Join(playersStr, "$"))

	if lc.cafe.GetRoomType() == cafe.CafeRoom {
		// Send every customer in location
		log.Debug("JOIN CUSTOMER DATA:")
		for _, cs := range lc.cafe.GetCustomers() {
			customerActionString := cs.ActionString()
			log.Debug("- SENT CUSTOMER DATA")

			if cs.GetAction() == customer.CUSTOMER_WALK_TO_CHAIR {
				customerActionString = cs.ActionStringToSpawnBack(customer.CUSTOMER_SIT_DOWN)
			}

			if cs.GetAction() != customer.CUSTOMER_LEAVE {
				lc.send(playerID, "nav", "-1", "0", cs.SpawnString())
				lc.send(playerID, "nac", "-1", "0", customerActionString)
			}
		}

		log.Debugf("CUSTOMER COUNT: %d ", len(lc.cafe.GetCustomers()))

		log.Debugf("WP COND %d %d", playerID, lc.cafe.GetPlayerID())

		// Start waiters when the owner joins if not yet stared
		if playerID == lc.cafe.GetPlayerID() && !lc.cafe.AgentCycleBinded {
			// Start waiters
			log.Debug("Waiters spawned and started")
			for i, w := range lc.cafe.Waiters {
				w.SetIsWorking(false)
				time.Sleep(10 * time.Millisecond)
				// Spawn waiters
				go func() {
					agents.SpawnWaiter(lc, w, i+1).Start()
				}()
			}

			lc.cafe.AgentCycleBinded = true
		} else if lc.cafe.AgentCycleBinded {
			// Respawn waiters
			for i, w := range lc.cafe.Waiters {

				waiterActionString := w.ActionString()
				if w.GetCurrentCounter() != nil {
					// NOTE: Always spawn the waiter to a counter. In the main loop, the waiter will be updated.
					waiterActionString = w.ActionStringToSpawnBack(waiter.MOVE_TO_COUNTER, w.GetCurrentCounter().GetPos())

				} else {
					// Fallback.
					waiterActionString = w.ActionStringToSpawnBack(waiter.INSERT, lc.cafe.GetPlayerStart())
				}
				lc.send(playerID, "nav", "-1", strconv.Itoa(i), w.SpawnString())
				lc.send(playerID, "nac", "-1", strconv.Itoa(i), waiterActionString)
			}
		}
	}

	lc.running = true
}

// Leaves the location and broadcasts leave to every one
func (lc *LoadedLocation) Leave(id int) {

	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Parse id
	idStr := strconv.Itoa(id)

	// Send leave to everyone
	lc.announce(id, "juq", "-1", "0", idStr)

	// Delete from occupants
	delete(lc.occupants, id)

	// If there are no players at the location
	// and the location is not marketplace
	// and the owner is not online
	if len(lc.occupants) == 0 && lc.cafe.GetID() > 0 && !lc.gm.UnsafeIsOnline(lc.cafe.GetID()) {
		log.Debugf("Unloading cafe %v...", lc.cafe.GetID())
		lc.gm.RemoveLocation(lc.cafe.ID)
	}

	// lc.running = false
}

func (lc *LoadedLocation) ClearReservedObjects() {
	for _, obj := range lc.reservedObjs {
		if obj.IsChair() {
			obj.SetDishID(-1)
			obj.SetDishStatus(0)
		}
	}
	lc.reservedObjs = []*object.Object{}
}

// Returns the first table and unreserves it
func (lc *LoadedLocation) GetDirtySpace() *object.Object {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	chairIndex := -1
	for i, o := range lc.reservedObjs {
		if o.IsChair() && o.GetDishStatus() == 3 {
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
		if o.GetPos().X == chair.GetPos().X+nr[0] && o.GetPos().Y == chair.GetPos().Y+nr[1] {
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
func (lc *LoadedLocation) GetReservedObject(pos simple.Position) *object.Object {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	for _, o := range lc.reservedObjs {
		if o.GetPos() == pos {
			return o
		}
	}

	return nil
}

func (lc *LoadedLocation) ReserveObject(obj *object.Object) bool {
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

func (lc *LoadedLocation) UnreserveObject(obj *object.Object) {
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

func (lc *LoadedLocation) Owner() (*player.Player, error) {
	c, err := lc.gm.GetClient(lc.cafe.GetID())
	if err != nil {
		return nil, err
	}
	return c.(*client.Client).Player, nil
}

// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
func (lc *LoadedLocation) send(id int, args ...string) {
	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-3), "to %v: %s", id, resp.Wrap())
	lc.occupants[id] <- resp
}

// | !!!  BEFORE USING THIS LOCK MUTEX  !!! |
// This will send a message to everyone in the loaded cafe
func (lc *LoadedLocation) broadcast(args ...string) {

	resp := responses.NewExtensionResponse(args...)
	log.Logf(log.Level(-1), "%s", resp.Wrap())
	// log.Debugf("Broadcasted for: %d clients", len(lc.occupants))
	for _, channel := range lc.occupants {
		channel <- resp
	}
}

// !!!  BEFORE USING THIS LOCK MUTEX  !!!
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
