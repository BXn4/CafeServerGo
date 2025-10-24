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
		10*time.Second, l) {
		l.Cafe().AddRating(-2)
		Leave(l, c) // Leaves sad :(
		return
	}

	if !l.IsRunning() {
		return
	}

	// Walk to chair and wait until arrives
	// println("Customer started to walking to the chair")
	CustomerDoAction(l, c,
		customer.CUSTOMER_WALK_TO_CHAIR,
		chair.GetPos(),
		time.Duration(distanceToChair-1)*550*time.Millisecond,
	)

	// println("Customer arrived to the chair")
	// println("Customer waiting 1sec before sit down")

	// Sit down
	CustomerDoAction(l, c,
		customer.CUSTOMER_SIT_DOWN,
		chair.GetPos(),
		1*time.Second,
	)

	// println("Customer sat down")

	// println("Customer started to waiting for food")

	// Wait for assigned waiter
	if WaitUntil(
		func() bool {
			return c.GetAssignedWaiter() != -1
		},
		25*time.Second, l) {
		l.Cafe().AddRating(-2)
		Leave(l, c) // Leaves sad :(
		l.UnreserveObject(table)
		l.UnreserveObject(chair)

		// println("Customer not got any food on time, customer left sad...")
		return
	}

	// Wait until waiter set food to the customer
	for c.GetDishID() == -1 {
		// Check if waiter abadoned customer
		if c.GetAssignedWaiter() == -1 {
			l.Cafe().AddRating(-2)
			Leave(l, c) // Leaves sad :( // We should make it to wait, then leave. TODO!
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
		if !l.TryStepSleep(100 * time.Millisecond) {
			return
		}
	}

	// println("Customer started to eating")

	//  Eat food
	CustomerDoAction(l, c, customer.CUSTOMER_EAT, c.GetPos(), 25*time.Second) // Customers eating for 25 sec. Footage: https://www.youtube.com/watch?v=pSX2kXIFxtE

	// After 10 sec set food to half eaten
	// I think we dont need to set the dish status. Just need to set when its finishes.
	// The game updates the dish status in visual to empty, when the customer finishes
	// chair.SetDishStatus(2) // Half eaten

	// println("Customer finished eating, leaving....")

	//  Add rewards to player
	player, err := l.Owner()
	if err != nil {
		log.Errorf("Cant find owner of cafe: %v", l.Cafe().GetPlayerID())
		return
	}

	// Get dish info
	dishInfo, err := utils.GetDish(c.GetDishID())
	if err != nil {
		log.Errorf("Cant find dish! %v\n", err)
		return
	}

	// Rewards
	player.AddCash(dishInfo.IncomePerServing)
	// The dish only gave XP when delivered to the counter
	// player.AddXP(dishInfo.XP)
	l.Cafe().AddRating(1)

	// Set plate dirty
	chair.SetDishStatus(3) // Dirty

	player.UpdateAchivementServingsCount(1)

	// c.DB.UpdateAchievement(player.ID, player.GetAchivements().String())

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
	// And instant cancel the sleeping it when the Café is not running
	if !l.TryStepSleep(5 * time.Second) {
		return
	}

	// Delete customer from customers
	l.Cafe().RemoveCustomer(c.GetID())
}

// Wait until function is true,
// returns true if timed out
func WaitUntil(condition func() bool, timeout time.Duration, l interfaces.CafeLocation) bool {
	startTime := time.Now()
	for !condition() {
		if time.Since(startTime) >= timeout {
			return true
		}
		if !l.TryStepSleep(10 * time.Millisecond) {
			return false
		}
	}
	return false
}

// 1. Sets action and pos,
// 2. Sends action response,
// 3. Waits for delay,
// 4. Stops if cafe is running,
func CustomerDoAction(l interfaces.CafeLocation, c *customer.Customer, action customer.CustomerAction, pos simple.Position, delay time.Duration) {

	// Set properties
	/* c.SetAction(action)
	   c.SetPos(pos)
	   I have moved these in the switch, because we set the customer action on sit down first, and theres 1 sec timer before the action,
	   so the waiters gave food when we set it to sit down action. We need to give them food after they have sat down */

	switch action {
	case customer.CUSTOMER_WALK_TO_CHAIR:
		// Set properties
		c.SetAction(action)
		c.SetPos(pos)

		// Send customer action
		l.Broadcast("nac", "-1", "0", c.ActionString())

		// Wait until customer walks in door
		if !l.TryStepSleep(delay) {
			return
		}
	case customer.CUSTOMER_SIT_DOWN:
		// Wait 1 sec before sit down (to avoid fast sit down at very close chair, because the distance is less than 1)
		if !l.TryStepSleep(delay) {
			return
		}

		// Set properties
		c.SetAction(action)
		c.SetPos(pos)

		l.Broadcast("nac", "-1", "0", c.ActionString())
	default:
		// Set properties
		c.SetAction(action)
		c.SetPos(pos)

		l.Broadcast("nac", "-1", "0", c.ActionString())
		if !l.TryStepSleep(delay) {
			return
		}
	}

}
