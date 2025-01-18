package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"cafego/internal/utils"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Spawns a custommer at the location
func SpawnCustomer(l interfaces.CafeLocation) *objects.Customer {

	rating := l.Cafe().Rating

	var spawnInterval int
	if rating < 150 {
		spawnInterval = rand.Intn(10) + 10
	} else if rating <= 150 && rating < 350 {
		spawnInterval = rand.Intn(3) + 5
	} else if rating <= 150 && rating < 350 {
		spawnInterval = rand.Intn(2) + 4
	} else {
		spawnInterval = rand.Intn(4) + 1
	}

	//spawnInterval = 2

	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(spawnInterval)*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return nil
		}
	}

	customer := &objects.Customer{
		ID:             l.GetUniqueCustomerID(),
		Avatar:         objects.NewRandomAvatar(),
		Pos:            []int{l.Cafe().PlayerStart[0], l.Cafe().PlayerStart[1]},
		Dish:           -1,
		Action:         objects.CUSTOMER_INSERT,
		IsThirsty:      false,
		AssignedWaiter: -2,
	}

	l.AddCustomer(customer)

	// Send customer info
	strID := strconv.Itoa(customer.ID)
	args := []string{
		strID,
		"0",  // NPC type (0: Customer)
		"0",  // Favourite = Waiters priority ???
		"-1", // DishID (unnecessary for waiters)
		utils.If(customer.IsThirsty, "1", "0"),
		customer.Avatar.String(),
	}
	// Send customer info + spawn
	l.Broadcast("nav", "-1", "0", strings.Join(args, "+"))
	l.Broadcast("nac", "-1", "0", strID+"+"+"0")

	return customer
}

// Does one iteration of the customer tasks
func IterateCustomer(l interfaces.CafeLocation, c *objects.Customer) {

	var table *objects.CafeObject
	var chair *objects.CafeObject
	var distanceToChair int

	// --- Wait until a eating space frees up ----------------------
	startTime := time.Now()
	for table == nil || chair == nil {

		if time.Since(startTime) >= 10*time.Second {
			Leave(l, c) // Leaves sad :(
			return
		}

		// Waiting for available space
		table, chair, distanceToChair = GetAvailableEatingSpace(l)
		time.Sleep(100 * time.Millisecond)
		if !l.IsRunning() {
			return
		}
	}

	// --- Walk to chair ---------------------------
	args := []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(objects.CUSTOMER_WALK_TO_CHAIR),
		strconv.Itoa(chair.Pos[0]),
		strconv.Itoa(chair.Pos[1]),
	}
	if !l.IsRunning() {
		l.UnreserveObject(table)
		l.UnreserveObject(chair)
		return
	} // We return if program is not running
	l.Broadcast("nac", "-1", "0", strings.Join(args, "+"))

	// Wait until walks to chair
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(distanceToChair-4)*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
	}
	if !l.IsRunning() {
		l.UnreserveObject(table)
		l.UnreserveObject(chair)
		return
	} // We return if program is not running

	// --- Sit down ---------------------------

	// Send
	c.Action = objects.CUSTOMER_SIT_DOWN
	args = []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(int(c.Action)),
		strconv.Itoa(chair.Pos[0]),
		strconv.Itoa(chair.Pos[1]),
	}
	if !l.IsRunning() {
		l.UnreserveObject(table)
		l.UnreserveObject(chair)
		return
	} // We return if program is not running
	l.Broadcast("nac", "-1", "0", strings.Join(args, "+"))

	// Set position to chair
	c.Pos[0] = chair.Pos[0]
	c.Pos[1] = chair.Pos[1]

	// Reset assigned waiter and food
	c.AssignedWaiter = -1
	c.Dish = -1

	// --- Wait for assigned waiter ----------------------
	startTime = time.Now()
	for c.AssignedWaiter == -1 {
		if !l.IsRunning() { // We return if program is not running
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
		if time.Since(startTime) >= 10*time.Second {
			Leave(l, c) // Leaves sad :(
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
		time.Sleep(1 * time.Second)
	}

	// --- Wait until food is placed ----------------------
	// startTime = time.Now()
	for c.Dish == -1 {
		if !l.IsRunning() { // We return if program is not running
			return
		}
		time.Sleep(1 * time.Second)
	}

	c.Action = objects.CUSTOMER_EAT

	// Wait for food so we dont eat the table xd
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < 3*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
	}

	// --- Eat food ---------------------------
	args = []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(int(c.Action)),
		strconv.Itoa(int(c.Pos[0])),
		strconv.Itoa(int(c.Pos[1])),
	}
	if !l.IsRunning() {
		l.UnreserveObject(table)
		l.UnreserveObject(chair)
		return
	} // We return if program is not running
	l.Broadcast(
		"nac", "-1", "0",
		strings.Join(args, "+"),
	)

	// Wait while checking for exit
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < 25*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			l.UnreserveObject(table)
			l.UnreserveObject(chair)
			return
		}
	}

	// --- Add rewards to player ------------
	player, err := l.Owner()
	if err != nil {
		log.Printf("Cant find owner!")
		return
	}

	dishInfo, err := utils.GetDish(c.Dish)
	if err != nil {
		log.Printf("Cant find dish! %v\n", c.Dish)
		return
	}
	player.Cash += dishInfo.IncomePerServing
	player.XP += dishInfo.XP

	// Dirty dishes
	chair.DishID = -2 // Dirty

	// --- Leave happy ------------------------
	LeaveComplete(l, c)
}

// Returns a chair and a table
// which are empty and approachable
func GetAvailableEatingSpace(l interfaces.CafeLocation) (*objects.CafeObject, *objects.CafeObject, int) {

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
			start := NewCafePoint(l.Cafe().PlayerStart, l)
			end := NewCafePoint(chair.Pos, l)
			_, distance, found := Path(start, end)
			if found {
				l.ReserveObject(chair)
				return table, chair, distance
			}
		}

	}
	return nil, nil, 0
}

func Leave(l interfaces.CafeLocation, c *objects.Customer) {

	// This will tell the waiter to not serve leaving customers
	c.Action = objects.CUSTOMER_LEAVE

	// Send leave
	args := []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(int(c.Action)),
	}
	if !l.IsRunning() { // We return if program is not running
		return
	}
	l.Broadcast(
		"nac", "-1", "0",
		strings.Join(args, "+"),
	)

	// Move to exit
	start := NewCafePoint(c.Pos, l)
	end := NewCafePoint(l.Cafe().PlayerStart, l)
	_, distance, _ := Path(start, end)

	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(distance)*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}

	// Delete customer from customers
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < 5*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}
	l.RemoveCustomer(c.ID)
}

func LeaveComplete(l interfaces.CafeLocation, c *objects.Customer) {

	// Set it do move
	c.Action = objects.CUSTOMER_LEAVE_COMPLETE
	c.Pos[0] = l.Cafe().PlayerStart[0]
	c.Pos[1] = l.Cafe().PlayerStart[1]

	// Send leave complete
	args := []string{
		strconv.Itoa(c.ID),
		strconv.Itoa(int(c.Action)),
	}
	if !l.IsRunning() {
		return
	} // We return if program is not running
	l.Broadcast(
		"nac", "-1", "0",
		strings.Join(args, "+"),
	)

	// Move to exit
	start := NewCafePoint(c.Pos, l)
	end := NewCafePoint(l.Cafe().PlayerStart, l)
	_, distance, _ := Path(start, end)
	time.Sleep(time.Duration(distance) * time.Second)

	// Delete customer from customers
	time.Sleep(5 * time.Second)
	l.RemoveCustomer(c.ID)
}
