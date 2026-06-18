package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
)

func init() {
	commands.RegisterCommand(requests.LOGIN,
		commands.CommandConfig{
			Name:       "RoomList",
			Identifier: responses.S2C_ROOMLIST,
			MinArgs:    1,
			MaxArgs:    2,
		},
		RoomListValidator,
		RoomList,
		nil,
	)
}

// rlu - RoomList
func RoomList(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	roomID := utils.If(req.Args[0] == responses.S2C_LOGIN, "1", "-1")

	c.SendExtensionResponse(responses.S2C_ROOMLIST, roomID, "1", "1", "20", "2", "Lobby")
	return nil
}

func RoomListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
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
