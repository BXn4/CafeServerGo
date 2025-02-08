package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/models/customer"
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"time"

	"github.com/charmbracelet/log"
)

// Spawns a custommer at the location
func SpawnCustomer(l interfaces.CafeLocation) *customer.Customer {

	// Create new random customer
	c := customer.NewRandomCustomer(
		l.GetUniqueCustomerID(),
		l.Cafe().GetPlayerStart(),
	)

	// Store customer data
	l.Cafe().AddCustomer(c)

	// Send customer action
	l.Broadcast("nav", "-1", "0", c.SpawnString())
	CustomerDoAction(l, c,
		customer.CUSTOMER_INSERT,
		l.Cafe().GetPlayerStart(),
		time.Second,
	)

	return c
}

// Does one iteration of the customer tasks
func IterateCustomer(l interfaces.CafeLocation, c *customer.Customer) {

	// Check if customer even exists
	if c == nil {
		return
	}

	// Declare varaibles
	var table *object.Object
	var chair *object.Object
	var distanceToChair int

	// Wait until a eating space frees up
	if WaitUntil(
		func() bool {
			table, chair, distanceToChair = GetAvailableEatingSpace(l)
			return chair != nil
		},
		10*time.Second,
	) {
		l.Cafe().AddRating(-2)
		Leave(l, c) // Leaves sad :(
		return
	}

	// Walk to chair and wait until arrives
	CustomerDoAction(l, c,
		customer.CUSTOMER_WALK_TO_CHAIR,
		chair.GetPos(),
		time.Duration(distanceToChair-2)*time.Second,
	)

	// Sit down
	CustomerDoAction(l, c, customer.CUSTOMER_SIT_DOWN, chair.GetPos(), 0)

	// Wait for assigned waiter
	if WaitUntil(
		func() bool {
			return c.GetAssignedWaiter() != -1
		},
		10*time.Second,
	) {
		l.Cafe().AddRating(-2)
		Leave(l, c) // Leaves sad :(
		l.UnreserveObject(table)
		l.UnreserveObject(chair)
		return
	}

	// Wait until food is placed
	for chair.GetDishStatus() != 1 {
		// Check if waiter abadoned customer
		if c.GetAssignedWaiter() == -1 {
			l.Cafe().AddRating(-2)
			Leave(l, c) // Leaves sad :(
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	//  Eat food
	CustomerDoAction(l, c, customer.CUSTOMER_EAT, c.GetPos(), 10*time.Second)

	// After 10 sec set food to half eaten
	chair.SetDishStatus(2) // Half eaten

	// Wait until customer finishes food
	time.Sleep(15 * time.Second)

	//  Add rewards to player
	player, err := l.Owner()
	if err != nil {
		log.Errorf("Cant find owner of cafe: %v", l.Cafe().GetPlayerID())
		return
	}

	// Get dish info
	dishInfo, err := utils.GetDish(chair.GetDishID())
	if err != nil {
		log.Errorf("Cant find dish! %v\n", chair.GetDishID())
		return
	}

	// Rewards
	player.AddCash(dishInfo.IncomePerServing)
	player.AddXP(dishInfo.XP)
	l.Cafe().AddRating(1)

	// Set plate dirty
	chair.SetDishStatus(3) // Dirty

	Leave(l, c)
}

// Returns a chair and a table
// which are empty and approachable
func GetAvailableEatingSpace(l interfaces.CafeLocation) (*object.Object, *object.Object, int) {

	spaces := l.Cafe().GetEatingSpaces()
	for table, chairs := range spaces {

		// If there are no connected chairs skip
		if len(chairs) == 0 {
			continue
		}

		// Try to reserve table
		if !l.ReserveObject(table) {
			continue
		}

		// Loop through all chairs and if approachable return them
		for _, chair := range chairs {
			start := NewCafePoint(l.Cafe().GetPlayerStart(), l.Cafe())
			end := NewCafePoint(chair.GetPos(), l.Cafe())
			_, distance, found := Path(start, end)
			if found {
				l.ReserveObject(chair)
				return table, chair, distance
			}
		}

	}
	return nil, nil, 0
}

func Leave(l interfaces.CafeLocation, c *customer.Customer) {

	// Calculate time to exit
	start := NewCafePoint(c.GetPos(), l.Cafe())
	end := NewCafePoint(l.Cafe().GetPlayerStart(), l.Cafe())
	_, distance, _ := Path(start, end)

	// Move out of cafe
	CustomerDoAction(l, c, customer.CUSTOMER_LEAVE, l.Cafe().GetPlayerStart(), time.Duration(distance)*time.Second)

	// Wait 5 sec to make sure client deleted customer
	time.Sleep(5 * time.Second)

	// Delete customer from customers
	l.Cafe().RemoveCustomer(c.GetID())
}

// Wait until function is true,
// returns true if timed out
func WaitUntil(condition func() bool, timeout time.Duration) bool {
	startTime := time.Now()
	for !condition() {
		if time.Since(startTime) >= timeout {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}

// 1. Sets action and pos,
// 2. Sends action response,
// 3. Waits for delay,
// 4. Stops if cafe is running,
func CustomerDoAction(l interfaces.CafeLocation, c *customer.Customer, action customer.CustomerAction, pos simple.Position, delay time.Duration) {

	// SEt properties
	c.SetAction(action)
	c.SetPos(pos)

	// Send customer action
	l.Broadcast("nac", "-1", "0", c.ActionString())

	// Wait until customer walks in door
	time.Sleep(delay)

	// Stop until cafe runs again
	for !l.IsRunning() {
		time.Sleep(10 * time.Millisecond)
	}
}
