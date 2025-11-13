package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"fmt"
	"time"
)

func CoopExtend(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	/*coopID, err := strconv.Atoi(req.Args[1]) not need
	if err != nil {
		return err
		} */

	if c.Player.GetActiveCoopID() == 0 {
		return nil
	}

	coop, err := c.DB.GetCoop(c.Player.CoopID)

	if err != nil {
		return fmt.Errorf("Could not get coop!")
	}

	if c.Player.GetGold() < balancing.BalancingConstants.CoopExpansionGold {
		c.SendExtensionResponse("cox", "-1", "4")
		return nil
	}

	c.Player.AddGold(-1)

	coop.AddExtend()

	coop.End = coop.End.Add(time.Duration(20) * time.Hour) // the time is fixed.
	// adds 20 hour to bronze, and we calculate the levels from the bronze time

	c.DB.SaveCoop(&coop)
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

	c.SendExtensionResponse("cox", "-1", "0%", coop.GetCoop().AsResponse())

	return nil
}
