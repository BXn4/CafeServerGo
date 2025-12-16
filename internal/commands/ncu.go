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
	RegisterCommand(requests.C2S_NPC_CUSTOMIZE,
		CommandConfig{
			Name:       "WaiterCustomize",
			Identifier: responses.S2C_NPC_CUSTOMIZE,
			MinArgs:    5,
			MaxArgs:    5,
		},
		WaiterCustomizeValidator,
		WaiterCustomize,
	)
}

func WaiterCustomize(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	selectedWaiter, _ := strconv.Atoi(req.Args[2])
	newPriority, _ := strconv.Atoi(req.Args[4])

	var w *waiter.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.GetID() == selectedWaiter {
			a := w.GetAvatar()
			a.Name = req.Args[3]
			w.SetAvatar(a)
			w.SetPriority(newPriority)
			break
		}
	}

	c.Location.Broadcast("ncu", "-1", "0",
		strconv.Itoa(w.GetID()),       // waiter id
		w.GetAvatar().Name,            // new waiter name
		strconv.Itoa(w.GetPriority()), // new waiter priority
	)

	c.DB.UpdateWaiters(c.Location.Cafe().ID, c.Location.Cafe().Waiters.String())

	return nil
}

func WaiterCustomizeValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	newPriority, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.GetLevel() < 3 {
		return "Player not yet unlocked the feature", NOT_DECLARED
	}

	var w *waiter.Waiter
	for _, w = range c.Location.Cafe().GetWaiters() {
		if w.GetID() == selectedWaiter {
			break
		}
	}

	if w == nil {
		return "Waiter not found!", NOT_DECLARED
	}

	if newPriority < 0 || newPriority > 100 {
		return "Priority is out of range!", NOT_DECLARED
	}

	if len(req.Args[3]) < 1 || len(req.Args) > 12 {
		return "New waiter name invalid!", ERROR_STAFF_RENAME_NAME_INVALID
	}

	return "Command ran without any errors.", NO_ERROR
}
