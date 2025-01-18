package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
	"time"

	"log"
)

func AgentCycle(l interfaces.CafeLocation) {

	log.Printf("---------------------------------\n")
	log.Printf("---------------------------------\n")
	log.Printf("-----  STARTED AGENT CYLCE  -----\n")
	log.Printf("---------------------------------\n")
	log.Printf("---------------------------------\n")

	l.ClearReservedObjects()

	if !SleepWhileRunning(l, 10*time.Second) {
		return
	}

	// Spawn waiters
	for i, w := range l.Cafe().Waiters {
		// Spawn waiter
		w.ID = i + 1
		SpawnWaiter(l, w)
	}

	// Count chairs
	var chairs []*objects.CafeObject
	for _, obj := range l.Cafe().Objects {
		if obj.IsChair() {
			chairs = append(chairs, obj)
		}
	}

	// Spawn customers
	go func() {
		for l.IsRunning() {
			log.Printf("CHAIRS LEN: %v\n", len(chairs))
			log.Printf("CUSTOMERS LEN: %v\n", len(l.Cafe().Customers))
			if len(l.Cafe().Customers) < len(chairs) {
				log.Printf("SPAWNED CUSTOMER!!!")
				go IterateCustomer(l, SpawnCustomer(l))
			}
		}
		l.Cafe().Customers = []*objects.Customer{}
	}()

	// IterateWaiters
	waiters := l.Cafe().Waiters
	for _, waiter := range waiters {
		go func() {
			for l.IsRunning() {
				IterateWaiter(l, waiter)
			}
			waiter.CurrentCounter = nil
			waiter.Dish = -1
		}()
	}

}

func SleepWhileRunning(l interfaces.CafeLocation, d time.Duration) bool {
	startTime := time.Now()
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for time.Since(startTime) < d {
		if !l.IsRunning() { // We return if program is not running
			return false
		}
		<-tick.C
	}
	return true
}
