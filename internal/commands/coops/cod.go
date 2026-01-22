package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"slices"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_COOP_DETAIL,
		commands.CommandConfig{
			Name:         "CoopDetail",
			Identifier:   responses.S2C_COOP_DETAIL,
			Description:  "Sending coop details",
			Args:         "{coop}",
			MinArgs:      3,
			MaxArgs:      3,
			FeatureLevel: 5,
		},
		CoopDetailValidator,
		CoopDetail,
		nil,
	)
}

// cod - CoopDetail
func CoopDetail(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	coopID, _ := strconv.Atoi(req.Args[2])

	if coopID == -2 {
		coopID = c.Player.GetCoopID()
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

	if slices.Contains(coop.Members, c.Player.GetID()) {
		c.SendExtensionResponse("cod", "-1", "0", strconv.Itoa(c.Player.GetID()), coop.GetCoop().AsResponse())
		if c.Player.GetCoopID() == coopID {
			if coop.FinishLevel != -1 {
				CoopFinish(&coop, gm)
				c.Player.SetCoopID(0)
				return nil
			}
		}

		return nil
	}

	c.SendExtensionResponse(cm.Identifier, "-1", "0%", coop.GetCoop().AsResponse())
	return nil
}

func CoopDetailValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet reached coops!", commands.NOT_DECLARED
	}

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if coopID == -2 {
		coopID = c.Player.GetCoopID()
	}

	if coopID != 0 {
		_, err = c.DB.GetCoop(coopID)
		if err != nil {
			return "Cant get coop detail from db!", commands.NOT_DECLARED
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
