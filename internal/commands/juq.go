package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// juq - S2C_JOIN_USERQUIT
func UserLeave(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	idStr := strconv.Itoa(c.Player.ID)

	// Send to people in the same location
	c.Cafe.Broadcast("juq", "-1", "0", idStr)

	return nil
}
