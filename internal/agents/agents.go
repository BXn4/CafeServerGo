package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/objects"
)

func AgentCycle(l interfaces.CafeLocation, isRunning bool) {

	l.ClearReservedObjects()

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
		for isRunning {
			if len(l.Cafe().Customers) < len(chairs) {
				go IterateCustomer(l, SpawnCustomer(l))
			}
		}
		l.Cafe().Customers = []*objects.Customer{}
	}()

	// IterateWaiters
	waiters := l.Cafe().Waiters
	for _, waiter := range waiters {
		go func() {
			for isRunning {
				IterateWaiter(l, waiter)
			}
		}()
	}

}
