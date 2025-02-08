package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

// Spawns a waiter at the location
func SpawnWaiter(l interfaces.CafeLocation, w *objects.Waiter) {

	w.IsWorking = true
	w.CurrentCounter = nil
	w.CurrentCustomer = nil

	// Set waiter starter position
	w.Pos = [2]int{
		l.Cafe().GetPlayerStart()[0],
		l.Cafe().GetPlayerStart()[1],
	}

	// Send waiter info
	args := []string{
		strconv.Itoa(w.ID),
		"1", // NPC type (1: Waiter)
		strconv.Itoa(w.Priority),
		"-1", // DishID (unnecessary for waiters)
		w.Avatar.String(w.Name),
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
	log.Debugf("WAITER SPAWNED: %v", w.ID)

}

// Does one iteration of the waiter tasks
func IterateWaiter(l interfaces.CafeLocation, w *objects.Waiter) {

	time.Sleep(1 * time.Second)

	job := rand.Intn(100) + 1 // 1-100

	if job > w.Priority {
		TakePlates(l, w)
	} else {
		ServeFood(l, w)
	}

}

func TakePlates(l interfaces.CafeLocation, w *objects.Waiter) {

	if !w.IsWorking { // We return if program is not running
		return
	}

	if w.CurrentCounter == nil {
		w.CurrentCounter, _ = GetRandomCounter(l.Cafe())
		if w.CurrentCounter == nil {
			return
		}
	}
	if !MoveWaiter(l, w, w.CurrentCounter.GetPos(), objects.MOVE_TO_COUNTER, 600*time.Millisecond) {
		return
	}

	// Get space with dirty plates
	space := l.GetDirtySpace()
	if space == nil {
		return
	}
	w.CurrentCounter = nil

	// Move to dirty plates
	if !MoveWaiter(l, w, space.GetPos(), objects.CLEAN, time.Duration(600)*time.Millisecond) {
		return
	}

	// Wait until waiter takes plates
	if !SleepWhileChecking(l, time.Second*5, &w.IsWorking) {
		return
	}

	// Bring back to counter
	for _, object := range l.Cafe().GetObjects() {
		if object.IsCounter() {
			if !MoveWaiter(l, w, space.GetPos(), objects.WAITER_MOVE_TO_COUNTER, time.Second) {
				return
			}

			w.CurrentCounter = object
			break
		}
	}

	// Wait and set the table clean
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
		counter, _ := GetRandomCounter(l.Cafe())
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
		if !MoveWaiter(l, w, w.CurrentCounter.GetPos(), objects.MOVE_TO_COUNTER, 600*time.Millisecond) {
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
	for _, c := range l.Cafe().GetCustomers() {
		if c.GetAction() == objects.CUSTOMER_SIT_DOWN && c.GetAssignedWaiter() == -1 {
			customer = c
			break
		}
	}

	// If every one has waiter return
	if customer == nil {
		return
	}

	// Assign itself as its waiter
	customer.SetAssignedWaiter(w.ID)

	// Take dish from counter prematurely so
	savedDish := w.CurrentCounter.GetDishID()
	w.CurrentCounter.SetDishAmount(w.CurrentCounter.GetDishAmount() - 1)
	if w.CurrentCounter.GetDishAmount() <= 0 {
		w.CurrentCounter.SetDishID(-1)
	}

	// --- Feed customer --------------------------
	if !MoveWaiter(l, w, customer.GetPos(), objects.FEED, 750*time.Millisecond) {
		return
	}

	// Set food to customer
	customer.SetDish(savedDish)

	// Move back to counter
	if !MoveWaiter(l, w, w.CurrentCounter.GetPos(), objects.MOVE_TO_COUNTER, 750*time.Millisecond) {
		return
	}
	w.CurrentCounter = nil
}

/*
Get a random counter,
that is reachable,
prioritizes counter with food,
return counter, distance
*/
func GetRandomCounter(cafe *objects.Cafe) (*objects.CafeObject, int) {

	var counters []*objects.CafeObject

	// Gather counters
	for _, object := range cafe.GetObjects() {

		// If object is not counter
		if !object.IsCounter() {
			continue
		}

		counters = append(counters, object)

		// Check if counter with food
		if object.GetDishID() >= 0 {

			// Check if blocked
			start := NewCafePoint([2]int(cafe.GetPlayerStart()), cafe)
			end := NewCafePoint(object.GetPos(), cafe)
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
		start := NewCafePoint([2]int(cafe.GetPlayerStart()), cafe)
		end := NewCafePoint(rc.GetPos(), cafe)
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
	start := NewCafePoint(w.Pos, l.Cafe())
	end := NewCafePoint(pos, l.Cafe())
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
