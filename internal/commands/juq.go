package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// juq - S2C_JOIN_USERQUIT
func UserLeave(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	idStr := strconv.Itoa(c.Player.ID)

	// Send to people in the same location
	c.Location.Broadcast("juq", "-1", "0", idStr)

	return nil
}
