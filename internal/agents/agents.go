package agents

import (
	"cafego/internal/interfaces"
	_ "cafego/internal/objects"
)

func AgentCycle(l interfaces.CafeLocation, isRunning bool) {
	// Spawn waiters
	for i, w := range l.Cafe().Waiters {
		// Spawn waiter
		w.ID = i + 1
		SpawnWaiter(l, w)
	}
	//var customers []*objects.Customer
	for isRunning {
		//TODO: Spawn customer every x sec

		// IterateWaiters
		IterateWaiters(l)

		//TODO: IterateCustomers
	}
}
