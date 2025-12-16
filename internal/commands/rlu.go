package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
)

func init() {
	RegisterCommand(requests.LOGIN,
		CommandConfig{
			Name:       "RoomList",
			Identifier: responses.S2C_ROOMLIST,
			MinArgs:    1,
			MaxArgs:    2,
		},
		RoomListValidator,
		RoomList,
	)
}

// rlu - RoomList
func RoomList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	roomID := utils.If(req.Args[0] == "lgn", "1", "-1")

	c.SendExtensionResponse("rlu", roomID, "1", "1", "20", "2", "Lobby")
	return nil
}

func RoomListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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
