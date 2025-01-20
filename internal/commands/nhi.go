package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
)

// nhi - C2S_NPC_HIRE
func HireWaiter(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	npcGender := utils.If(req.Args[3] == "1", objects.Girl, objects.Boy)
	avatar := objects.NewRandomAvatar()
	avatar.Gender = npcGender
	avatar.Name = req.Args[2]
	newWaiter := &objects.Waiter{
		ID:       len(c.Location.Cafe().Waiters) + 1,
		Name:     req.Args[2],
		Priority: 50,
		Avatar:   avatar,
	}
	c.Location.Cafe().Waiters = append(c.Location.Cafe().Waiters, newWaiter)

	// Spawn
	agents.SpawnWaiter(c.Location, newWaiter)

	// Start waiter cylce
	go func() {
		for c.Location.IsRunning() {
			agents.IterateWaiter(c.Location, newWaiter)
		}
		newWaiter.CurrentCounter = nil
		newWaiter.Dish = -1
	}()

	c.SendExtensionResponse("nhi", "0", "0", req.Args[2], req.Args[3])
	return nil
}
