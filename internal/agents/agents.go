package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"time"
)

func AgentCycle(l interfaces.CafeLocation) {

	l.ClearReservedObjects()

	// Spawn waiters
	for i, w := range l.Cafe().GetWaiters() {
		// Spawn waiter
		w.IsWorking = true
		w.CurrentCounter = nil
		w.CurrentCustomer = nil
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

func SleepWhileChecking(l interfaces.CafeLocation, d time.Duration, isRunning *bool) bool {
	startTime := time.Now()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for time.Since(startTime) < d {
		if !*isRunning { // We return if program is not running
			return false
		}
		<-tick.C
	}
	return true
}
