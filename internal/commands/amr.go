package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_ALLOW_MAIL_REQUESTS,
		CommandConfig{
			Name:       "AllowMailRequests",
			Identifier: responses.S2C_ALLOW_MAIL_REQUESTS,
			MinArgs:    3,
			MaxArgs:    3,
			IsBool:     true,
		},
		AllowEmailsValidator,
		AllowEmails,
	)
}

// abr - AllowEmails
func AllowEmails(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	c.Player.AllowEmails = utils.If(req.Args[2] == "1", true, false)

	c.SendExtensionResponse("amr", "-1", req.Args[2], strconv.Itoa(c.Player.ID))

	c.DB.UpdateAllowEmails(c.Player.ID, c.Player.AllowEmails)
	return nil
}

func AllowEmailsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	return "Command ran without any errors.", NO_ERROR
}
