package agents

import (
	"cafego/internal/interfaces"
	"cafego/internal/models/cafe"
	"cafego/internal/models/object"
	"math"
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
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

	l.Cafe().AgentCycleBinded = true

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		nextSpawn := time.Now().UTC().Add(GetSpawnInterval(l))

		for range ticker.C {

			if !l.Cafe().AgentCycleBinded {
				println("Stopping agent cycle: not binded")
				return
			}

			if !*l.GetIsRunning() {
				println("Agent cycle paused")
				continue
			}

			if time.Now().UTC().After(nextSpawn) {
				nextSpawn = time.Now().UTC().Add(GetSpawnInterval(l))

				go Spawn(l, NewCustomer(l))
			}

			// println("TICK!")
		}
	}()
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

func GetSpawnInterval(l interfaces.CafeLocation) time.Duration {
	maxSpawn := 30.0
	minSpawn := 2.0
	rating := float64(l.Cafe().GetRating())
	expansion := float64(l.Cafe().ExpansionID)
	ratingFactor := math.Min(rating/1000.0, 10.0)
	expansionFactor := math.Min(expansion/8.0, 1.0)
	progress := ratingFactor*0.6 + expansionFactor*0.4
	spawnBase := maxSpawn - progress*(maxSpawn-minSpawn)
	variation := 0.8 + rand.Float64()*0.4
	nextSpawn := time.Duration(spawnBase*variation) * time.Second
	log.Debugf("NPC spawn interval: %s", nextSpawn.String())
	return nextSpawn
}
