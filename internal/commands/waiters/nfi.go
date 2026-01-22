package waiters

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/waiter"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_NPC_FIRE,
		commands.CommandConfig{
			Name:        "WaiterFire",
			Identifier:  responses.S2C_NPC_FIRE,
			Description: "Fire an waiter from the Café",
			Args:        "{waiterID}",
			MinArgs:     3,
			MaxArgs:     3,
		},
		WaiterFireValidator,
		WaiterFire,
		WaiterFireDBSaver,
	)
}

// nfi
func WaiterFire(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {

	id, _ := strconv.Atoi(req.Args[2])

	c.Location.Cafe().RemoveWaiter(id)

	c.Player.UpdateAchivementFireNPC()

	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())

	c.SendExtensionResponse(cm.Identifier, "0", "0", req.Args[2])

	return nil
}

func WaiterFireValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	selectedWaiter, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if c.Player.GetLevel() < 3 {
		return "Player not yet unlocked the feature", commands.ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	var w *waiter.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.GetID() == selectedWaiter {
			break
		}
	}

	if w == nil {
		return "Waiter not found!", commands.ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	if len(c.Location.Cafe().GetWaiters()) == 1 {
		return "Cant fire the last waiter!", commands.ERROR_STAFF_FIRE_IMPOSSIBLE
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func WaiterFireDBSaver(c *client.Client) error {
	c.DB.UpdateWaiters(c.Location.Cafe().GetID(), c.Location.Cafe().GetWaiters().String())

	return nil
}
