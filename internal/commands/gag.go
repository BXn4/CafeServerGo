package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

func DailyGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// TODO: Generate a set of items every 24 hours (best to do it with daily login)

	c.SendExtensionResponse("gag", "-1", "0", "")
	return nil
}
