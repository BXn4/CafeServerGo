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

func (lc *LoadedLocation) Cafe() *cafe.Cafe {
	if lc.cafe == nil {
		log.Errorf("Attempted to access nil cafe in location")
		return nil
	}
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

	if c.(*client.Client).Player.GetIsSeekingJob() {
		c.(*client.Client).Player.SetIsSeekingJob(false)
	} // need to clear it

	// Set position of player
	c.(*client.Client).Player.SetPos(lc.cafe.GetPlayerStart())

	// Send everyone the joined player
	lc.announce(playerID, "juj", "-1", "0", c.(*client.Client).Player.String())

	// Add to players in cafe
	lc.occupants[playerID] = channel

	playersStr := []string{}

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

	julArgs := append([]string{"jul", "-1", "0"}, playersStr...)

	// Send it to joined player
	lc.send(playerID, julArgs...)

	// If the room is cafe
	if lc.cafe.GetRoomType() == cafe.CafeRoom {
		switch lc.running {
		case true:
			//println("LOCATION IS RUNNING")
			// Respawn waiters
			for _, w := range lc.cafe.GetWaiters() {
				waiterActionString := w.ActionString()
				if w.GetCurrentCounter() != nil {
					// NOTE: Always spawn the waiter to a counter. In the main loop, the waiter will be updated.
					waiterActionString = w.ActionStringToSpawnBack(waiter.MOVE_TO_COUNTER, w.GetCurrentCounter().GetPos())

				} else {
					// Fallback.
					waiterActionString = w.ActionStringToSpawnBack(waiter.INSERT, lc.cafe.GetPlayerStart())
				}
				lc.send(playerID, "nav", "-1", "0", w.SpawnString())
				lc.send(playerID, "nac", "-1", "0", waiterActionString)
			}

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
		case false:
			// println("LOCATION IS NOT RUNNING")
			if playerID == lc.Cafe().GetOwnerID() {
				p, err := lc.Owner()
				if err != nil {
					log.Errorf("Cant find owner of cafe: %v", lc.Cafe().GetOwnerID())
				}

				if p == nil {
					return
				}

				if p.GetIsTutorialCompleted() {
					go agents.FillEmptyCafe(lc)
					go agents.StartAgentCycles(lc)
				}

				for i, w := range lc.cafe.GetWaiters() {
					w.SetIsWorking(false)
					// Spawn waiters
					go func() {
						agents.SpawnWaiter(lc, w, i+1).Start()
					}()
				}

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
	if lc.cafe != nil && len(lc.occupants) == 0 && lc.cafe.GetID() > 0 && !lc.gm.UnsafeIsOnline(lc.cafe.GetOwnerID()) {
		log.Debugf("Unloading cafe %v...", lc.cafe.GetID())
		lc.gm.RemoveLocation(lc.cafe.GetID())
	}

	// lc.running = false
}

func (lc *LoadedLocation) ClearReservedObjects() {
	for _, obj := range lc.reservedObjs {
		if obj.IsChair() {
			obj.SetDishID(-1)
			obj.SetDishStatus(0)
			obj.SetOccupied(false)
		}
	}
	lc.reservedObjs = []*object.Object{}
}

// Returns the first table
func (lc *LoadedLocation) GetDirtySpace() (*object.Object, *object.Object) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	chairIndex := -1
	for i, o := range lc.reservedObjs {
		if o.IsChair() && o.GetDishStatus() == 3 {
			chairIndex = i
		}
	}

	if chairIndex == -1 {
		return nil, nil
	}

	chair := lc.reservedObjs[chairIndex]

	// Get associated table
	for _, o := range lc.reservedObjs {
		if !o.IsTable() {
			continue
		}
		nr := chair.GetNormalizedRotation()
		if o.GetPos().X == chair.GetPos().X+nr[0] && o.GetPos().Y == chair.GetPos().Y+nr[1] {
			table := o
			return table, chair
		}
	}

	return nil, nil
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

	obj.SetOccupied(true)

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

	obj.SetOccupied(false)

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

func (lc *LoadedLocation) TryStepSleep(d time.Duration) bool {
	step := 100 * time.Millisecond
	elapsed := time.Duration(0)
	for elapsed < d {
		if !*lc.GetIsRunning() {
			return false
		}

		time.Sleep(step)
		elapsed += step
	}

	return true
}
