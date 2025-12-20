package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/customer"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_EDITOR_MODE,
		CommandConfig{
			Name:       "EditorMode",
			Identifier: responses.S2C_EDITOR_BUY_OBJECT,
			MinArgs:    3,
			MaxArgs:    3,
			IsBool:     true,
		},
		EditorModeValidator,
		EditorMode,
	)
}

// edit - C2S_EDITOR_MODE
func EditorMode(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	status, _ := strconv.Atoi(req.Args[2])

	switch status {
	case 0:
		c.Location.SetRunning(true)
		c.Location.ClearReservedObjects()
		cafe := c.Location.Cafe()
		if cafe == nil {
			return fmt.Errorf("Location cafe is nil for client %d", c.ClientID)
		}
		c.Player.Position = cafe.GetPlayerStart()

		cafe.Customers = make(map[int]*customer.Customer)

		go agents.FillEmptyCafe(c.Location)

		for i, w := range cafe.Waiters {
			w.SetIsWorking(false)
			// Spawn waiters
			go func() {
				agents.SpawnWaiter(c.Location, w, i+1).Start()
			}()
		}
	case 1:
		c.Location.Cafe().SetCustomers(nil)
		for _, w := range c.Location.Cafe().Waiters {
			w.StopWorking()
		}
		c.Location.SetRunning(false)
	}

	c.SendExtensionResponse("edi", "-1", "0",
		strconv.Itoa(status),
		strconv.Itoa(c.Player.Position.X),
		strconv.Itoa(c.Player.Position.Y),
		"", // <- Objects
	)

	return nil
}

func EditorModeValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if cm.IsBool {
		if req.Args[2] != "0" && req.Args[2] != "1" {
			return fmt.Sprintf("Invalid args for boolean: %v", req.Args[2]), INVALID_ARGS
		}
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	if !c.Location.IsRunning() && req.Args[2] != "0" {
		return "Cant enter editor mode again!", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
