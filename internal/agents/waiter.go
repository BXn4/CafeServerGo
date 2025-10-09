package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/models/customer"
	"cafego/internal/models/simple"
	"cafego/internal/models/waiter"
	"math/rand"
	"time"
)

type TaskFunction func() TaskFunction

type WaiterAgent struct {
	l  interfaces.CafeLocation
	w  *waiter.Waiter
	fn TaskFunction
}

func SpawnWaiter(l interfaces.CafeLocation, w *waiter.Waiter, id int) *WaiterAgent {

	if w.IsWorking() {
		return nil
	}

	wa := &WaiterAgent{l: l, w: w}

	// Reset waiter properties
	wa.w.SetIsWorking(true)
	wa.w.SetID(id)
	wa.w.SetPos(l.Cafe().GetPlayerStart())
	wa.w.SetCurrentCounter(nil)
	wa.w.SetCurrentCustomer(nil)

	// Send waiter spawn response
	wa.l.Broadcast("nav", "-1", "0", wa.w.SpawnString())
	/*if wa.sleep(time.Second * 1) {
	return nil
	} */

	println("SPAWNED WAITER: ", id, w.GetAvatar().Name)

	// Send action response
	wa.doAction(waiter.INSERT, wa.w.GetPos(), 1*time.Second)

	return wa
}

// Starts main waiter loop
// This helps take full controll of waiter actions
func (wa *WaiterAgent) Start() {
	// If waiter doesent exits return
	if wa == nil {
		return
	}

	println("STARTED WAITER TASK LOOP: ", wa.w.GetID())
	// Start waiter task loop
	go func() {
		defer wa.Stop()
		wa.fn = wa.getAndMoveToCounter
		for wa.fn != nil {
			wa.fn = wa.fn()
			if !wa.w.IsWorking() {
				break
			}
		}
	}()
}

// Stops waiter loop
func (wa *WaiterAgent) Stop() {
	wa.w.SetIsWorking(false)
	if wa.w.GetCurrentCustomer() != nil {
		wa.w.GetCurrentCustomer().SetAssignedWaiter(-1)
	}
}

// --- Private methods -----------------

func (wa *WaiterAgent) getAndMoveToCounter() TaskFunction {

	// if wa.sleep(time.Second * 2) {
	// 	return nil
	// }

	if wa.w.GetCurrentCounter() == nil || wa.w.GetCurrentCounter().GetDishID() == -1 {

		// Get random counter
		counter, _ := GetRandomCounter(wa.l.Cafe())
		if counter == nil {
			return nil
		}

		// If counter has food change it
		if wa.w.GetCurrentCounter() == nil || (counter.GetDishID() != -1 && wa.w.GetCurrentCounter().GetDishID() == -1) {
			wa.w.SetCurrentCounter(counter)
		}

		// Move to counter
		if wa.move(wa.w.GetCurrentCounter().GetPos(), waiter.MOVE_TO_COUNTER) {
			return nil
		}

	}
	return wa.selectJob
}

func (wa *WaiterAgent) selectJob() TaskFunction {

	// Roll work (like in original)
	job := rand.Intn(100) + 1 // 1-100

	// Do work based on priority
	if job > wa.w.GetPriority() {
		return wa.takePlates
	} else {
		return wa.serveFood
	}
}

func (wa *WaiterAgent) takePlates() TaskFunction {

	// If we dont have a counter return
	if wa.w.GetCurrentCounter() == nil {
		return wa.getAndMoveToCounter
	}

	// Get space with dirty plates
	space := wa.l.GetDirtySpace()
	if space == nil {
		//println("CANT FIND DIRTY SPACE")
		return wa.getAndMoveToCounter
	}

	println("MOVE TO CLEAN")

	// Set space clean
	// Moved here, because we need to remove the dirty dishes from the table for the new client who joins the Café, because after it, we cant and it remains.
	// Also I added SetFree and GetIsFree for the chair, so customers now checking this
	space.SetDishID(-1)
	space.SetDishStatus(0)

	// Move to dirty plates
	if wa.move(space.GetPos(), waiter.CLEAN) {
		return nil
	}

	println("WAIT UNTIL CLEANED")
	//Wait until waiter takes plates
	if wa.sleep(time.Second * 2) {
		return nil
	}

	space.SetOccupied(false)

	wa.w.SetCurrentCounter(nil)

	return wa.getAndMoveToCounter
}

func (wa *WaiterAgent) serveFood() TaskFunction {

	// If we dont have a counter return
	if wa.w.GetCurrentCounter() == nil {
		return wa.getAndMoveToCounter
	}

	// If counter is empty
	if wa.w.GetCurrentCounter().GetDishID() == -1 {
		return wa.getAndMoveToCounter
	}

	// Get sitting customer without waiter
	var cu *customer.Customer
	for _, c := range wa.l.Cafe().GetCustomers() {
		if c.GetAction() == customer.CUSTOMER_SIT_DOWN && c.GetAssignedWaiter() == -1 {
			cu = c
			break
		}
	}

	// If every one has waiter return
	if cu == nil {
		return wa.getAndMoveToCounter
	}

	// Assign itself as its waiter
	cu.SetAssignedWaiter(wa.w.GetID())

	// Take dish from counter prematurely so there will be one there
	savedDish := wa.w.GetCurrentCounter().GetDishID()
	wa.w.GetCurrentCounter().SetDishAmount(wa.w.GetCurrentCounter().GetDishAmount() - 1)
	if wa.w.GetCurrentCounter().GetDishAmount() <= 0 {
		wa.w.GetCurrentCounter().SetDishID(-1)
	}

	chair := wa.l.Cafe().GetObjectByPos(cu.GetPos())

	// Set food to customer
	// Set it earlier, because we need to send the chair dish id to the visitor, and if we not sending it on time, the customer will eat noting, and gives negative rating in visual
	chair.SetDishID(savedDish)
	chair.SetDishStatus(1)

	// Move to customer and feed customer
	if wa.move(cu.GetPos(), waiter.FEED) {
		return nil
	}

	if wa.sleep(time.Second * 1) {
		return nil
	}

	// Waiter arrived, food delivered
	cu.SetDishID(savedDish)

	// Check
	// if wa.w.GetCurrentCounter() == nil {
	// 	return nil
	// }

	// Move back to counter
	if wa.move(wa.w.GetCurrentCounter().GetPos(), waiter.MOVE_TO_COUNTER) {
		return nil
	}

	// Reset counter
	wa.w.SetCurrentCounter(nil)

	return wa.getAndMoveToCounter
}

// --- Helper methods ---------------------

// This searches for a path and moves it if possible
func (wa *WaiterAgent) move(pos simple.Position, action waiter.Action) bool {

	// Get distance to location
	start := NewCafePoint(wa.w.GetPos(), wa.l.Cafe())
	end := NewCafePoint(pos, wa.l.Cafe())
	path, distance, found := Path(start, end)
	if !found {
		return false
	}

	if distance <= 1 {
		return false
	}

	// Set pos
	wa.w.SetPos(path[1].Pos())

	println("distance: ", distance)

	// Set waiter pos
	wa.doAction(action, pos, time.Duration(distance)*450*time.Millisecond)

	println("Waiter arrived to the destination pos")

	return false
}

// Does action if it gets interupted returns true
func (wa *WaiterAgent) doAction(action waiter.Action, pos simple.Position, delay time.Duration) bool {

	// Set properties
	wa.w.SetAction(action)
	savedPos := wa.w.GetPos()
	wa.w.SetPos(pos)

	// Send task response
	wa.l.Broadcast("nac", "-1", "0", wa.w.ActionString())

	// Load last pos
	wa.w.SetPos(savedPos)

	// Wait until task ends
	return wa.sleep(delay)
}

// Almost like time.Sleep but if it gets interupted returns true
func (wa *WaiterAgent) sleep(delay time.Duration) bool {

	// Escape if time expires
	for {
		// Wait spawn time
		ticker := time.NewTicker(10 * time.Millisecond) // Tick every 100 ms
		defer ticker.Stop()
		expire := time.After(delay)
		for {
			select {
			case <-ticker.C:
				if !wa.w.IsWorking() {
					return true
				}
			case <-expire:
				return false
			}
		}
	}
}
