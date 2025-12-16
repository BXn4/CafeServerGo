package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/waiter"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_NPC_FIRE,
		CommandConfig{
			Name:       "WaiterFire",
			Identifier: responses.S2C_NPC_FIRE,
			MinArgs:    3,
			MaxArgs:    3,
		},
		WaiterFireValidator,
		WaiterFire,
	)
}

// nfi
func WaiterFire(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	id, _ := strconv.Atoi(req.Args[2])

	c.Location.Cafe().RemoveWaiter(id)

	c.Player.UpdateAchivementFireNPC()

	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

	c.SendExtensionResponse("nfi", "0", "0", req.Args[2])

	c.DB.UpdateWaiters(c.Location.Cafe().ID, c.Location.Cafe().Waiters.String())

	return nil
}

func WaiterFireValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	selectedWaiter, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.GetLevel() < 3 {
		return "Player not yet unlocked the feature", ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	var w *waiter.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.GetID() == selectedWaiter {
			break
		}
	}

	if w == nil {
		return "Waiter not found!", ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	if len(c.Location.Cafe().GetWaiters()) == 1 {
		return "Cant fire the last waiter!", ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	return "Command ran without any errors.", NO_ERROR
}
