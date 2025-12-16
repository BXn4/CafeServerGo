package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/models/balancing"
	"cafego/internal/models/waiter"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
)

func init() {
	RegisterCommand(requests.C2S_NPC_HIRE,
		CommandConfig{
			Name:       "WaiterHire",
			Identifier: responses.S2C_NPC_HIRE,
			MinArgs:    4,
			MaxArgs:    4,
		},
		WaiterHireValidator,
		WaiterHire,
	)
}

// min level 3

// nhi - C2S_NPC_HIRE
func WaiterHire(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// Pay price
	c.Player.AddGold(-balancing.BalancingConstants.StaffPrice)

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

	newWaiter := waiter.NewWaiter(newID, 50, avatar, false)

	c.Location.Cafe().AddWaiter(newWaiter)

	// Start waiter agent cycle
	go func() {
		// Spawn waiter and start waiter cylce
		agents.SpawnWaiter(c.Location, newWaiter, newID).Start()
	}()

	println(c.Location.Cafe().Waiters.String())

	c.DB.UpdateWaiters(c.Location.Cafe().ID, c.Location.Cafe().Waiters.String())

	c.SendExtensionResponse("nhi", "0", "0", req.Args[2], req.Args[3])
	return nil
}

func WaiterHireValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	if c.Player.GetLevel() < 3 {
		return "Player not yet unlocked the feature", ERROR_STAFF_HIRE_IMPOSSIBLE
	}

	if len(c.Location.Cafe().GetWaiters()) >= utils.GetLevelWaitersLimit(c.Player.GetLevel()) {
		return "Cant hire more waiters, because the level limit!", ERROR_STAFF_HIRE_IMPOSSIBLE
	}

	if c.Player.GetGold() < balancing.BalancingConstants.StaffPrice {
		return "Player not have enough money", NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", NO_ERROR
}
