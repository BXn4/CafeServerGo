package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
)

// col - CoopLeave
func CoopLeave(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	if !c.Player.IsInCoop() {
		return nil
	}

	coop, err := c.DB.GetCoop(c.Player.GetActiveCoopID())
	if err != nil {
		return fmt.Errorf("Cannot get coop!")
	}

	coop.Leave(c.Player.ID)

	if len(coop.Members) == 0 {
		c.DB.DeleteCoop(coop.ID)
	} else {
		c.DB.SaveCoop(&coop) // coop rebuilds when leaving
	}

	/* playersString := make([]string, len(coop.Members))

	for i, playerID := range coop.Members {
		player, err := c.DB.GetPlayer(playerID)

		if err != nil {
			playersString[i] = ""
		} else {
			playersString[i] = strconv.Itoa(player.ID) + "+" + strconv.Itoa(player.XP) + "+" + player.Avatar.String()
		}
	}

	coop.SetPlayersString(playersString) */

	c.Player.SetActiveCoopID(0)
	c.DB.UpdateCoopID(c.Player.ID, c.Player.GetActiveCoopID())
	c.SendExtensionResponse("col", "-1", "0")

	return nil
}
