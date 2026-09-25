/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_COOP_ACTIVELIST,
		commands.CommandConfig{
			Name:         "ActiveCoops",
			Identifier:   responses.S2C_COOP_ACTIVELIST,
			Description:  "Sending the active coops list to the client",
			Args:         "{coopList}#",
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 5,
		},
		CoopActiveListValidator,
		CoopActiveList,
		nil,
	)
}

// coa - CoopActiveList
func CoopActiveList(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	var args []string
	if c.Player.GetIsInCoop() {
		coop, _ := c.DB.GetCoop(c.Player.GetCoopID())
		args = append(args, coop.GetCoop().AsActiveListResponse())
	}

	for _, playerID := range c.Player.GetFriends() {
		coop, err := c.DB.GetCoopByHost(playerID)
		if err == nil {
			if coop.ID != c.Player.GetCoopID() {
				if coop.GetIsActive() {
					args = append(args, coop.AsActiveListResponse())
				}
			}
		}
	}

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strings.Join(args, "#"))
	return nil
}

func CoopActiveListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet reached coops!", commands.NOT_DECLARED
	}

	if c.Player.GetIsInCoop() {
		_, err := c.DB.GetCoop(c.Player.GetCoopID())
		if err != nil {
			return "Cannot get player active coop!", commands.NOT_DECLARED
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
