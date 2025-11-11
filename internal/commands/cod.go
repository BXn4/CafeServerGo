package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"slices"
	"strconv"
)

// cod - CoopDetail
func CoopDetail(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	coopID, err := strconv.Atoi(req.Args[2])

	if err != nil {
		return fmt.Errorf("Cannot parse coopID to int!")
	}

	if coopID == -2 {
		coopID = c.Player.GetActiveCoopID()
	}

	if coopID == 0 {
		c.SendExtensionResponse("cod", "-1", "0")
		return nil
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return fmt.Errorf("Cannot get coop!")
	}

	/* if(_loc2_)
	 * this._activeCoop = _loc7_;
	 * if loc2 is given (memberid) it sets the coop as active to the client
	 */

	if slices.Contains(coop.Members, c.Player.ID) {
		c.SendExtensionResponse("cod", "-1", "0", strconv.Itoa(c.Player.ID), coop.GetCoop().AsResponse())
		if c.Player.GetActiveCoopID() == coopID {
			if coop.FinishLevel != -1 {
				CoopFinish(&coop, gm)
				c.Player.SetActiveCoopID(0)
				c.DB.UpdateCoopID(c.Player.ID, c.Player.GetActiveCoopID())
				return nil
			}
		}

		return nil
	}

	c.SendExtensionResponse("cod", "-1", "0%", coop.GetCoop().AsResponse())
	return nil
}
