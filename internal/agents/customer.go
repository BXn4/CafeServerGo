/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/models/balancing"
	"cafego/internal/models/customer"
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"math"
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
)

type EatingSpace struct {
	Chair    *object.Object
	Table    *object.Object
	Distance int
}

func FillEmptyCafe(l interfaces.CafeLocation) {
	// Number of customers based on:
	// The chairs and the rating

	// Spawn random eating customers, and random waiting customers
	// If theres no food on the counters -> waiting
	// If theres food on the counters -> 1/2 eat or wait

	// The eating not works, need to set the chairs before the objects was sent
	numberOfChairs := 0
	rating := float64(l.Cafe().GetRating())
	expansion := float64(l.Cafe().GetExpansionID())
	ratingFactor := math.Min(rating/1000.0, 10.0)
	expansionFactor := math.Min(expansion/8.0, 1.0)

	spaces := l.Cafe().GetEatingSpaces()
	for _, chairs := range spaces {

		// If there are no connected chairs skip
		if len(chairs) == 0 {
			continue
		}

		// Loop through all chairs and if approachable count them
		for _, chair := range chairs {
			start := NewCafePoint(l.Cafe().GetPlayerStart(), l.Cafe())
			end := NewCafePoint(chair.GetPos(), l.Cafe())
			_, _, found := Path(start, end)
			if found {
				numberOfChairs++
			}
		}
	}

	customersToSpawn := int(float64(numberOfChairs) * ratingFactor * expansionFactor)

	if numberOfChairs > 0 && customersToSpawn == 0 {
		// spawn 1 OR 2 50 % force
		customersToSpawn = int(math.Min(float64(numberOfChairs), float64(rand.Intn(3))))
	}

	println("Customers to spawn: %d", customersToSpawn)

	for _ = range customersToSpawn {
		// customerCanEat := false
		c := NewCustomer(l)

		/* counter, _ := GetRandomCounter(l.Cafe())

		// if its not going to be empty
		if counter.GetDishAmount() > 0 {
			customerCanEat = true
		}

		customerIsEating := rand.Intn(2) == 1 */

		table, chair, distanceToChair := GetAvailableEatingSpace(l)
		es := EatingSpace{Chair: chair, Table: table, Distance: distanceToChair}

		/*if customerIsEating && customerCanEat {
			c.SetDishID(counter.GetDishID())
			es.Chair.SetDishID(counter.GetDishID())
			counter.AddDishAmount(-1)
		}*/

		l.Broadcast("nav", "-1", "0", c.SpawnString())

		/*if customerIsEating && customerCanEat {
			eatingSec := rand.Intn(25-5) + 5

			c.SetAction(customer.CUSTOMER_EAT)
			c.SetPos(es.Chair.GetPos())

			l.Broadcast("nac", "-1", "0", c.ActionString())

			go Eat(l, c, time.Duration(eatingSec), es)

		} else { */
		waitingSec := rand.Intn(25-7) + 7
		c.SetAction(customer.CUSTOMER_SIT_DOWN)
		c.SetPos(es.Chair.GetPos())

		l.Broadcast("nac", "-1", "0", c.ActionString())

		go WaitForFood(l, c, time.Duration(waitingSec), es)
		// }
	}
}

// Add the customer to cafe customer list, returns a new customer
func NewCustomer(l interfaces.CafeLocation) *customer.Customer {

	if !*l.GetIsRunning() {
		return nil
	}

	// Create new random customer
	c := customer.NewRandomCustomer(
		l.GetUniqueCustomerID(),
		l.Cafe().GetPlayerStart(),
	)

	// Store customer data
	l.Cafe().AddCustomer(c)

	return c
}

// Spawns a custommer at the location
func Spawn(l interfaces.CafeLocation, c *customer.Customer) {

	if c == nil {
		return
	}

	// Send customer action
	l.Broadcast("nav", "-1", "0", c.SpawnString())
	if CustomerDoAction(l, c,
		customer.CUSTOMER_INSERT,
		l.Cafe().GetPlayerStart(),
		time.Second,
	) {
		log.Debug("Customer task cancelled!")
		return
	}

	WaitForChair(l, c)
}

func WaitForChair(l interfaces.CafeLocation, c *customer.Customer) {
	if WaitUntil(
		func() bool {
			table, chair, distanceToChair := GetAvailableEatingSpace(l)
			es := EatingSpace{Chair: chair, Table: table, Distance: distanceToChair}
			if es.Chair != nil && es.Table != nil && es.Distance != -1 {
				WalkToChair(l, c, es)
				return true
			}
			return false
		},
		10*time.Second, l) {
		Leave(l, c, nil)
		return
	}
}

func WalkToChair(l interfaces.CafeLocation, c *customer.Customer, es EatingSpace) {
	println("Customer started to walking to the chair")
	if CustomerDoAction(l, c,
		customer.CUSTOMER_WALK_TO_CHAIR,
		es.Chair.GetPos(),
		time.Duration(es.Distance-1)*550*time.Millisecond,
	) {
		log.Debug("Customer task cancelled!")
		return
	}

	println("Customer arrived to the chair")
	SitDown(l, c, es)
}

func SitDown(l interfaces.CafeLocation, c *customer.Customer, es EatingSpace) {
	println("Customer waiting 1sec before sit down")

	if CustomerDoAction(l, c,
		customer.CUSTOMER_SIT_DOWN,
		es.Chair.GetPos(),
		1*time.Second,
	) {
		log.Debug("Customer task cancelled!")
		return
	}

	println("Customer sat down")
	WaitForFood(l, c, time.Duration(25), es)
}

func WaitForFood(l interfaces.CafeLocation, c *customer.Customer, duration time.Duration, es EatingSpace) {
	println("Customer started to waiting for food")

	if WaitUntil(
		func() bool {
			if c.GetAssignedWaiter() != -1 {
				WaitForWaiterArrive(l, c, es)
				return true
			}
			return false
		},
		duration*time.Second, l) {
		println("Customer not got any food on time, customer left sad...")
		Leave(l, c, &es) // Leaves sad :(
		return
	}
}

func WaitForWaiterArrive(l interfaces.CafeLocation, c *customer.Customer, es EatingSpace) {
	println("Customer started to waiting for waiter to arrive")
	for c.GetDishID() == -1 {
		// Check if waiter abadoned customer
		if c.GetAssignedWaiter() == -1 {
			println("Waiter abadoned the customer!")
			// 10, because I dont want to make it to wait again 25s
			WaitForFood(l, c, time.Duration(10), es)
			return
		}
		if !l.TryStepSleep(100 * time.Millisecond) {
			return
		}
	}

	Eat(l, c, time.Duration(25), es)
	// Customers eating for 25 sec. Footage: https://www.youtube.com/watch?v=pSX2kXIFxtE
}

func Eat(l interfaces.CafeLocation, c *customer.Customer, duration time.Duration, es EatingSpace) {
	println("Customer started to eat")
	switch {
	case duration >= 15:
		es.Chair.SetDishStatus(1)
	case duration >= 5:
		es.Chair.SetDishStatus(2)
	case duration <= 5:
		es.Chair.SetDishStatus(3)
	}

	if CustomerDoAction(l, c, customer.CUSTOMER_EAT, c.GetPos(), duration*time.Second) {
		log.Debug("Customer task cancelled!")
		return
	}

	println("Customer ate the food!")
	Leave(l, c, &es)
}

func Leave(l interfaces.CafeLocation, c *customer.Customer, es *EatingSpace) {
	// Leave sad, not got any dish
	if c.GetDishID() == -1 {
		l.Cafe().AddRating(balancing.BalancingConstants.RatingGuestUnhappy)

		// If the customer had a chair, but not got food
		// So the waiter not need to clean it
		if es != nil {
			l.UnreserveObject(es.Chair)
			l.UnreserveObject(es.Table)
		}
	} else {
		// Customer got food!
		player, err := l.Owner()
		if err != nil {
			log.Errorf("Cant find owner of cafe: %v", l.Cafe().GetOwnerID())
			return
		}

		dishInfo, err := utils.GetDish(c.GetDishID())
		if err != nil {
			log.Errorf("Cant find dish! %v\n", err)
			return
		}

		player.AddCash(dishInfo.IncomePerServing)
		// The dish only gave XP when delivered to the counter
		// player.AddXP(dishInfo.XP)
		l.Cafe().AddRating(balancing.BalancingConstants.RatingGuestHappy)

		es.Chair.SetDishStatus(3) // Dirty

		player.UpdateAchivementServingsCount(1)
	}

	// Calculate time to exit
	start := NewCafePoint(c.GetPos(), l.Cafe())
	end := NewCafePoint(l.Cafe().GetPlayerStart(), l.Cafe())
	_, distance, _ := Path(start, end)

	// Move out of cafe
	if CustomerDoAction(l, c, customer.CUSTOMER_LEAVE, l.Cafe().GetPlayerStart(), time.Duration(distance)*time.Second) {
		// Wait 5 sec to make sure client deleted customer
		// And instant cancel the sleeping it when the Café is not running
		log.Debug("Customer task cancelled!")
		return
	}

	if !l.TryStepSleep(5 * time.Second) {
		log.Debug("Customer task cancelled!")
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

// 1. Sets action and pos,
// 2. Sends action response,
// 3. Waits for delay,
// 4. Stops if cafe is running,
// 5. Return bool false if its cancelled, true when mpt
func CustomerDoAction(l interfaces.CafeLocation, c *customer.Customer, action customer.CustomerAction, pos simple.Position, delay time.Duration) bool {

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
			println("RETURN TRUE")
			return true
		}
	case customer.CUSTOMER_SIT_DOWN:
		// Wait 1 sec before sit down (to avoid fast sit down at very close chair, because the distance is less than 1)
		if !l.TryStepSleep(delay) {
			return true
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
			return true
		}
	}

	return false
}
