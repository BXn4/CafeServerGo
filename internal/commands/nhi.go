package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/models/waiter"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
)

// nhi - C2S_NPC_HIRE
func HireWaiter(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Checl if we have enough money
	if c.Player.GetGold() < 2 {
		return nil
	}

	// Pay price
	c.Player.AddGold(-2)

	npcGender := utils.If(req.Args[3] == "1", avatar.Girl, avatar.Boy)
	avatar := avatar.NewRandomAvatar()
	avatar.Gender = npcGender
	avatar.Name = req.Args[2]

	// find new waiter id
	newID := 1
	changed := true
	for changed {
		changed = false
		for _, w := range c.Location.Cafe().GetWaiters() {
			if newID == w.GetID() {
				newID++
				changed = true
			}
		}
	}

	// create new waiter
	newWaiter := waiter.NewWaiter(newID, 50, avatar, false)

	c.Location.Cafe().AddWaiter(newWaiter)

	// Start waiter agent cycle
	go func() {
		// Spawn waiter and start waiter cylce
		agents.SpawnWaiter(c.Location, newWaiter, newID).Start()
	}()

	c.DB.UpdateWaiters(c.Location.Cafe().ID, c.Location.Cafe().Waiters.String())

	c.SendExtensionResponse("nhi", "0", "0", req.Args[2], req.Args[3])
	return nil
}
