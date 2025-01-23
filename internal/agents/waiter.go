package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Spawns a waiter at the location
func SpawnWaiter(l interfaces.CafeLocation, w *objects.Waiter) {

	// Set waiter starter position
	w.Pos = [2]int{
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
	if !w.IsWorking {
		return
	}
	l.Broadcast("nav", "-1", "0", strings.Join(args, "+"))

	// Spawn waiter
	l.Broadcast("nac", "-1", "0", args[0], "0")

	// Wait to go in door
	if !SleepWhileChecking(l, 1*time.Second, &w.IsWorking) {
		return
	}

	// --- Spawn waiter ----------
	log.Printf("WAITER SPAWNED: %v\n", w.ID)

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

	if !w.IsWorking { // We return if program is not running
		return
	}

	// Get space with dirty plates
	space := l.GetDirtySpace()
	if space == nil {
		return
	}
	w.CurrentCounter = nil

	// Move to dirty plates
	if !MoveWaiter(l, w, space.GetPos(), objects.CLEAN, time.Duration(500)*time.Millisecond) {
		return
	}

	// Bring back to counter
	for _, object := range l.Cafe().Objects {
		if object.IsCounter() {

			if !MoveWaiter(l, w, space.GetPos(), objects.WAITER_MOVE_TO_COUNTER, time.Second) {
				return
			}

			w.CurrentCounter = object
			break
		}
	}

	// Wait for response and set the table clean
	if !SleepWhileChecking(l, time.Second, &w.IsWorking) {
		return
	}
	space.SetDishID(-1)

	// Wait until it puts back to counter
	if !SleepWhileChecking(l, time.Second*3, &w.IsWorking) {
		return
	}

	return
}

func ServeFood(l interfaces.CafeLocation, w *objects.Waiter) {

	if !w.IsWorking { // We return if program is not running
		return
	}

	//var distance int
	if w.CurrentCounter == nil || w.CurrentCounter.GetDishID() == -1 {
		// --- Get random counter (prioritizes counter with food) -------------------------------
		counter, _ := GetRandomCounter(l)
		if counter == nil {
			return
		}

		// If counter has food change it
		if w.CurrentCounter == nil {
			w.CurrentCounter = counter
		} else if counter.GetDishID() != -1 && w.CurrentCounter.GetDishID() == -1 {
			w.CurrentCounter = counter
		}

		// --- Move to counter ------------------------------------
		if !MoveWaiter(l, w, w.CurrentCounter.GetPos(), objects.MOVE_TO_COUNTER, 500*time.Millisecond) {
			return
		}
	}

	// If could not found counter (DONT TOUCH IT !!! IT WORKS!!!)
	if w.CurrentCounter == nil {
		return
	}
	// If counter is empty
	if w.CurrentCounter.GetDishID() == -1 {
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
	savedDish := w.CurrentCounter.GetDishID()
	w.CurrentCounter.SetDishAmount(w.CurrentCounter.GetDishAmount() - 1)
	if w.CurrentCounter.GetDishAmount() <= 0 {
		w.CurrentCounter.SetDishID(-1)
	}

	// --- Feed customer --------------------------
	if !MoveWaiter(l, w, customer.Pos, objects.FEED, 750*time.Millisecond) {
		return
	}

	// Set food to customer
	customer.Dish = savedDish

	// Move back to counter
	if !MoveWaiter(l, w, w.CurrentCounter.GetPos(), objects.MOVE_TO_COUNTER, 750*time.Millisecond) {
		return
	}
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
		if object.GetDishID() >= 0 {

			// Check if blocked
			start := NewCafePoint(l.Cafe().PlayerStart, l)
			end := NewCafePoint(object.GetPos(), l)
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
		start := NewCafePoint(l.Cafe().PlayerStart, l)
		end := NewCafePoint(rc.GetPos(), l)
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
func MoveWaiter(l interfaces.CafeLocation, w *objects.Waiter, pos [2]int, action objects.Action, duration time.Duration) bool {

	// Get length of path
	start := NewCafePoint(w.Pos, l)
	end := NewCafePoint(pos, l)
	path, distance, found := Path(start, end)
	if !found {
		return false
	}

	// Send move msg
	waiterPos := path[1]
	w.Pos[0] = waiterPos.x
	w.Pos[1] = waiterPos.y
	args := []string{
		strconv.Itoa(w.ID),
		strconv.Itoa(int(action)),
		strconv.Itoa(pos[0]),
		strconv.Itoa(pos[1]),
	}

	if !w.IsWorking {
		return false
	}
	l.Broadcast("nac", "-1", "0", strings.Join(args, "+"))

	if !SleepWhileChecking(l, time.Duration(distance)*duration, &w.IsWorking) {
		return false
	}

	return found
}
