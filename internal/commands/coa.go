package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strings"
)

func init() {
	RegisterCommand(requests.C2S_COOP_ACTIVELIST,
		CommandConfig{
			Name:       "ActiveCoops",
			Identifier: responses.S2C_COOP_ACTIVELIST,
			MinArgs:    0,
			MaxArgs:    0,
		},
		CoopActiveListValidator,
		CoopActiveList,
	)
}

// coa - CoopActiveList
func CoopActiveList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var args []string
	if c.Player.IsInCoop() {
		coop, _ := c.DB.GetCoop(c.Player.GetActiveCoopID())
		args = append(args, coop.GetCoop().AsActiveListResponse())
	}

	for _, playerID := range c.Player.Friends {
		coop, err := c.DB.GetCoopByHost(playerID)
		if err == nil {
			if coop.ID != c.Player.GetActiveCoopID() {
				if coop.GetIsActive() {
					args = append(args, coop.AsActiveListResponse())
				}
			}
		}
	}

	c.SendExtensionResponse("coa", "-1", "0", strings.Join(args, "#"))
	return nil
}

func CoopActiveListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if c.Player.IsInCoop() {
		_, err := c.DB.GetCoop(c.Player.GetActiveCoopID())
		if err != nil {
			return "Cannot get player active coop!", NOT_DECLARED
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
