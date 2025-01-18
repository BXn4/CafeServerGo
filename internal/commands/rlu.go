package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
)

// rlu - RoomList
func RoomList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	roomID := utils.If(req.Args[0] == "lgn", "1", "-1")

	c.SendExtensionResponse("rlu", roomID, "1", "1", "20", "2", "Lobby")
	return nil
}
