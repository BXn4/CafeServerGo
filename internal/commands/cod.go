package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"slices"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_COOP_DETAIL,
		CommandConfig{
			Name:       "CoopDetail",
			Identifier: responses.S2C_COOP_DETAIL,
			MinArgs:    3,
			MaxArgs:    3,
		},
		CoopDetailValidator,
		CoopDetail,
	)
}

// cod - CoopDetail
func CoopDetail(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	coopID, _ := strconv.Atoi(req.Args[2])

	if coopID == -2 {
		coopID = c.Player.GetActiveCoopID()
	}

	if coopID == 0 {
		c.SendExtensionResponse("cod", "-1", "0")
		return nil
	}

	coop, _ := c.DB.GetCoop(coopID)

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

func CoopDetailValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if c.Player.GetLevel() < 5 {
		return "Player not yet reached coops!", CONVERT_ERROR
	}

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if coopID == -2 {
		coopID = c.Player.GetActiveCoopID()
	}

	if coopID != 0 {
		_, err = c.DB.GetCoop(coopID)
		if err != nil {
			return "Cant get coop detail from db!", NOT_DECLARED
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
