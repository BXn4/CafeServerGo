package settings

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
)

func init() {
	commands.RegisterCommand(requests.C2S_CHANGE_PASSWORD,
		commands.CommandConfig{
			Name:        "ChangePassword",
			Identifier:  responses.S2C_CHANGE_PASSWORD,
			Description: "Change password",
			Args:        "{}",
			MinArgs:     5,
			MaxArgs:     5,
		},
		ChangePasswordValidator,
		ChangePassword,
		nil,
	)
}

// lcp - ChangePassword
func ChangePassword(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	oldPasswd, newPasswd := req.Args[2], req.Args[3]

	status, err := c.DB.ChangePassword(c.ID(), oldPasswd, newPasswd)
	if err != nil {
		log.Warnf("Error changing password: %v", err)
	}

	c.SendExtensionResponse(cm.Identifier, "1", strconv.Itoa(status))
	if status == 0 {
		c.Disconnect()
	}
	return nil
}

func ChangePasswordValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
