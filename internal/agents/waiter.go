package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Spawns a waiter at the location
func SpawnWaiter(l interfaces.CafeLocation, w *objects.Waiter) {

	// --- Spawn waiter ----------
	println("WAITER SPAWNED: ", w.ID)

	// Set waiter starter position
	w.Pos = []int{
		l.Cafe().PlayerStart[0],
		l.Cafe().PlayerStart[1],
	}

	// Send waiter info
	args := []string{
		strconv.Itoa(w.ID),
		"1", // NPC type (1: Waiter)
		strconv.Itoa(w.Priority),
		"-1", // DishID (unnecessary for waiters)
		w.Avatar.String(),
	}
	if !l.IsRunning() {
		return
	} // We return if program is not running
	l.Broadcast("nav", "-1", "0", strings.Join(args, "+"))

	// Spawn waiter
	l.Broadcast("nac", "-1", "0", args[0], "0")

}

// Does one iteration of the waiter tasks
func IterateWaiter(l interfaces.CafeLocation, w *objects.Waiter) {

	time.Sleep(1 * time.Second)
	if w.Priority == objects.BOTH {
		ServeFood(l, w)
		TakePlates(l, w)
	} else if w.Priority == objects.SERVING {
		// - Serve food while there are customers
		// - if there are no customers take the plates
	} else if w.Priority == objects.CLEANING {
		// - Take plates while there are plates
		// - if there are no plates serve customers
	}
}

func TakePlates(l interfaces.CafeLocation, w *objects.Waiter) {

	if !l.IsRunning() { // We return if program is not running
		return
	}

	// Get space with dirty plates
	space := l.GetDirtySpace()
	if space == nil {
		return
	}
	w.CurrentCounter = nil

	// Move to dirty plates
	_, distance := MoveWaiter(l, w, space.Pos[0], space.Pos[1], objects.CLEAN)
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(500*distance)*time.Millisecond; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}

	// Bring back to counter
	for _, object := range l.Cafe().Objects {
		if object.IsCounter() {
			_, distance := MoveWaiter(l, w, space.Pos[0], space.Pos[1], objects.WAITER_MOVE_TO_COUNTER)
			for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(distance)*time.Second; <-tick.C {
				if !l.IsRunning() { // We return if program is not running
					return
				}
			}
			w.CurrentCounter = object
			break
		}
	}

	// Wait for response and set the table clean
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < 1*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}
	space.DishID = -1

	// Wait until it puts back to counter
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < 3*time.Second; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}

	return
}

func ServeFood(l interfaces.CafeLocation, w *objects.Waiter) {

	if !l.IsRunning() { // We return if program is not running
		return
	}

	var distance int
	if w.CurrentCounter == nil || w.CurrentCounter.DishID == -1 {
		// --- Get random counter -------------------------------
		counter, distance := GetRandomCounter(l)
		if counter == nil {
			return
		}

		// If counter has food change it
		if w.CurrentCounter == nil {
			w.CurrentCounter = counter
		} else if counter.DishID != -1 && w.CurrentCounter.DishID == -1 {
			w.CurrentCounter = counter
		}

		// --- Move to counter ------------------------------------
		_, distance = MoveWaiter(l, w, w.CurrentCounter.Pos[0], w.CurrentCounter.Pos[1], objects.MOVE_TO_COUNTER)
		//time.Sleep(time.Duration(750*distance) * time.Millisecond)
		for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(500*distance)*time.Millisecond; <-tick.C {
			if !l.IsRunning() { // We return if program is not running
				return
			}
		}
	}

	//
	if w.CurrentCounter.DishID == -1 {
		return
	}

	// Get sitting customer without waiter
	var customer *objects.Customer
	for _, c := range l.Cafe().Customers {
		if c.Action == objects.CUSTOMER_SIT_DOWN && c.AssignedWaiter == -1 {
			customer = c
			break
		}
	}

	// If every one has waiter return
	if customer == nil {
		return
	}

	// Assign itself as its waiter
	customer.AssignedWaiter = w.ID

	// Take dish from counter prematurely so
	savedDish := w.CurrentCounter.DishID
	w.CurrentCounter.DishAmount -= 1
	if w.CurrentCounter.DishAmount <= 0 {
		w.CurrentCounter.DishID = -1
	}

	// --- Feed customer --------------------------
	_, distance = MoveWaiter(l, w, customer.Pos[0], customer.Pos[1], objects.FEED)
	for start, tick := time.Now(), time.NewTicker(100*time.Millisecond); time.Since(start) < time.Duration(750*distance)*time.Millisecond; <-tick.C {
		if !l.IsRunning() { // We return if program is not running
			return
		}
	}

	// Set food to customer
	customer.Dish = savedDish

	// Move back to counter
	_, distance = MoveWaiter(l, w, w.CurrentCounter.Pos[0], w.CurrentCounter.Pos[1], objects.MOVE_TO_COUNTER)

	w.CurrentCounter = nil

}

// Get a random counter that has food and it is reachable from the start location
func GetRandomCounter(l interfaces.CafeLocation) (*objects.CafeObject, int) {

	var counters []*objects.CafeObject

	for _, object := range l.Cafe().Objects {

		if !object.IsCounter() {
			continue
		}

		counters = append(counters, object)

		// Check if counter with food
		if object.DishID >= 0 {

			// Check if blocked
			start := &CafePoint{x: l.Cafe().PlayerStart[0], y: l.Cafe().PlayerStart[1], l: l}
			end := &CafePoint{x: object.Pos[0], y: object.Pos[1], l: l}
			_, distance, found := Path(start, end)

			// If found path there return it
			if found {
				return object, distance
			}
		}
	}

	// check if th
	for len(counters) != 0 {

		i := rand.Intn(len(counters))
		rc := counters[i] // random counter

		// Search path
		start := &CafePoint{x: l.Cafe().PlayerStart[0], y: l.Cafe().PlayerStart[1], l: l}
		end := &CafePoint{x: rc.Pos[0], y: rc.Pos[1], l: l}
		_, distance, found := Path(start, end)

		// If found path return
		if found {
			return rc, distance
		}
		counters = append(counters[:i], counters[i+1:]...)
	}

	return nil, -1
}

// This searches for a path and moves it if possible
// returns if found path and length of the found path
func MoveWaiter(l interfaces.CafeLocation, w *objects.Waiter, x int, y int, action objects.Action) (bool, int) {

	// Get length of path
	start := &CafePoint{x: w.Pos[0], y: w.Pos[1], l: l}
	end := &CafePoint{x: x, y: y, l: l}
	path, distance, found := Path(start, end)

	if !found {
		return false, -1
	}

	// Send move msg
	waiterPos := path[1]
	w.Pos[0] = waiterPos.x
	w.Pos[1] = waiterPos.y
	args := []string{
		strconv.Itoa(w.ID),
		strconv.Itoa(int(action)),
		strconv.Itoa(x),
		strconv.Itoa(y),
	}
	if !l.IsRunning() {
		return false, -1
	} // We return if program is not running
	l.Broadcast("nac", "-1", "0", strings.Join(args, "+"))

	return found, distance
}
