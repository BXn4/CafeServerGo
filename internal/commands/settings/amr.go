/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package settings

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_ALLOW_MAIL_REQUESTS,
		commands.CommandConfig{
			Name:        "AllowMailRequests",
			Identifier:  responses.S2C_ALLOW_MAIL_REQUESTS,
			Description: "Enable/Disable emails.",
			Args:        "{0/1}",
			MinArgs:     3,
			MaxArgs:     3,
			IsBool:      true,
		},
		AllowEmailsValidator,
		AllowEmails,
		AllowEmailsDBSaver,
	)
}

// abr - AllowEmails
func AllowEmails(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	c.Player.SetAllowEmails(utils.If(req.Args[2] == "1", true, false))
	c.SendExtensionResponse(cm.Identifier, "-1", req.Args[2], strconv.Itoa(c.Player.GetID()))
	return nil
}

func AllowEmailsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if cm.IsBool {
		if req.Args[2] != "0" && req.Args[2] != "1" {
			return fmt.Sprintf("Invalid args for boolean: %v", req.Args[2]), commands.INVALID_ARGS
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func AllowEmailsDBSaver(c *client.Client) error {
	c.DB.UpdateAllowEmails(c.Player.GetID(), c.Player.GetAllowEmails())

	return nil
}
