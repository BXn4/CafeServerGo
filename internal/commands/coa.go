package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strings"
)

// coa - CoopActiveList
func CoopActiveList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var args []string
	if c.Player.IsInCoop() {
		coop, err := c.DB.GetCoop(c.Player.GetActiveCoopID())

		if err != nil {
			return fmt.Errorf("Cannot get coop")
		}

		args = append(args, coop.GetCoop().AsActiveListResponse())
	}

	for _, playerID := range c.Player.Friends {
		coop, err := c.DB.GetCoopByHost(playerID)
		if err == nil {
			if coop.ID != c.Player.GetActiveCoopID() {
				if coop.GetIsActive() {
					args = append(args, coop.AsActiveListResponse())
				}
			}
		}
	}

	c.SendExtensionResponse("coa", "-1", "0", strings.Join(args, "#"))
	return nil
}
