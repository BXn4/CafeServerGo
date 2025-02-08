package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/models/cafe"
	"cafego/internal/models/object"
	"math/rand"
	"time"
)

func StartAgentCycles(l interfaces.CafeLocation) {

	// Clean tables and chairs
	for _, obj := range l.Cafe().GetObjects() {
		if obj.IsTable() || obj.IsChair() {
			obj.SetDishID(-1)
		}
	}

	// Empty reserved objects
	l.ClearReservedObjects()

	// Spawn customers
	go func() {
		for l != nil {
			// Calcualte customer spawn time
			rating := l.Cafe().GetRating()
			var spawnInterval int
			if rating < 150 {
				spawnInterval = rand.Intn(10) + 10
			} else if rating <= 150 && rating < 350 {
				spawnInterval = rand.Intn(3) + 5
			} else if rating <= 350 && rating < 500 {
				spawnInterval = rand.Intn(2) + 4
			} else {
				spawnInterval = rand.Intn(4) + 1
			}
			time.Sleep(time.Duration(spawnInterval) * time.Second)

			// Stop if we are not in area
			for !l.IsRunning() {
				if l == nil {
					return
				}
			}

			go IterateCustomer(l, SpawnCustomer(l))

		}
	}()

	// IterateWaiters
	// waiters := l.Cafe().GetWaiters()
	// for i, w := range waiters {
	// 	// Main cycle
	// 	go func() {
	// 		SpawnWaiter(l, w, i+1).Start()
	// 	}()
	// }
}

/*
Get a random counter,
that is reachable,
prioritizes counter with food,
return counter, distance
*/
func GetRandomCounter(c *cafe.Cafe) (*object.Object, int) {
	var counters []*object.Object

	// Gather counters and check for ones with food
	for _, obj := range c.GetObjects() {
		if !obj.IsCounter() {
			continue
		}

		counters = append(counters, obj)

		if obj.GetDishID() >= 0 {
			start := NewCafePoint(c.GetPlayerStart(), c)
			end := NewCafePoint(obj.GetPos(), c)
			_, distance, found := Path(start, end)
			if found {
				return obj, distance
			}
		}
	}

	// Try random counters if no counter with food is found
	for len(counters) > 0 {
		i := rand.Intn(len(counters))
		rc := counters[i]
		start := NewCafePoint(c.GetPlayerStart(), c)
		end := NewCafePoint(rc.GetPos(), c)
		_, distance, found := Path(start, end)

		if found {
			return rc, distance
		}
		counters = append(counters[:i], counters[i+1:]...)
	}

	return nil, -1
}
