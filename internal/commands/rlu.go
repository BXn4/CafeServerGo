package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// rlu - RoomList
func RoomList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	var roomID string
	if req.Args[0] == "lgn" {
		roomID = "1"
	} else {
		roomID = "-1"
	}

	c.SendExtensionResponse("rlu", roomID, "1", "1", "20", "2", "Lobby")
	return nil
}
