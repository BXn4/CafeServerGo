package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
)

func init() {
	RegisterCommand(requests.C2S_CHANGE_PASSWORD,
		CommandConfig{
			Name:       "ChangePassword",
			Identifier: responses.S2C_CHANGE_PASSWORD,
			MinArgs:    5,
			MaxArgs:    5,
		},
		ChangePasswordValidator,
		ChangePassword,
	)
}

// lcp - ChangePassword
func ChangePassword(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	oldPasswd, newPasswd := req.Args[2], req.Args[3]

	status, err := c.DB.ChangePassword(c.ID(), oldPasswd, newPasswd)
	if err != nil {
		log.Warnf("Error changing password: %v", err)
	}

	c.SendExtensionResponse("lcp", "1", strconv.Itoa(status))
	if status == 0 {
		c.Disconnect()
	}
	return nil
}

func ChangePasswordValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
