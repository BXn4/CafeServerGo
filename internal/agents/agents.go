package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"time"
)

func AgentCycle(l interfaces.CafeLocation) {

	// Clean tables and chairs
	for _, obj := range l.Cafe().GetObjects() {
		if obj.IsTable() || obj.IsChair() {
			obj.SetDishID(-1)
		}
	}

	//
	l.ClearReservedObjects()

	// Spawn waiters
	for i, w := range l.Cafe().GetWaiters() {
		// Spawn waiter
		w.ID = i + 1
		SpawnWaiter(l, w)
	}

	// Count chairs
	var chairs []*objects.CafeObject
	for _, obj := range l.Cafe().GetObjects() {
		if obj.IsChair() {
			chairs = append(chairs, obj)
		}
	}

	// Spawn customers
	go func() {
		for l.IsRunning() {
			if len(l.Cafe().GetCustomers()) < len(chairs) {
				go IterateCustomer(l, SpawnCustomer(l))
			}
		}
		l.Cafe().SetCustomers([]*objects.Customer{})
	}()

	// IterateWaiters
	waiters := l.Cafe().GetWaiters()
	for _, waiter := range waiters {

		go func() {
			for waiter.IsWorking {
				IterateWaiter(l, waiter)
			}
			waiter.CurrentCounter = nil
			waiter.Dish = -1
		}()
	}

}

// This is a sleep wrapper that checks every 100ms the paramater boolean
// and returns false if it stopped being true
func SleepWhileChecking(l interfaces.CafeLocation, d time.Duration, isRunning *bool) bool {
	startTime := time.Now()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for time.Since(startTime) < d {
		// If program is not running
		if !*isRunning {
			return false
		}
		<-tick.C
	}
	return true
}
