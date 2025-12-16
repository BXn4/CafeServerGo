package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"time"
)

func init() {
	RegisterCommand(requests.C2S_COOP_EXTEND,
		CommandConfig{
			Name:       "CoopExtend",
			Identifier: responses.S2C_COOP_EXTEND,
			MinArgs:    0,
			MaxArgs:    0,
		},
		CoopExtendValidator,
		CoopExtend,
	)
}

func CoopExtend(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	/*coopID, err := strconv.Atoi(req.Args[1]) not need
	if err != nil {
		return err
		} */

	coop, _ := c.DB.GetCoop(c.Player.CoopID)
	c.Player.AddGold(-1)
	coop.AddExtend()
	coop.End = coop.End.Add(time.Duration(20) * time.Hour) // the time is fixed.
	// adds 20 hour to bronze, and we calculate the levels from the bronze time

	c.DB.SaveCoop(&coop)
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

	c.SendExtensionResponse("cox", "-1", "0%", coop.GetCoop().AsResponse())

	return nil
}

func CoopExtendValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	if c.Player.GetActiveCoopID() == 0 {
		return "Cant extend coop, because the player not in coop!", NOT_DECLARED
	}

	_, err := c.DB.GetCoop(c.Player.CoopID)
	if err != nil {
		return "Cant get coop detail from db!", NOT_DECLARED
	}

	if c.Player.GetGold() < balancing.BalancingConstants.CoopExpansionGold {
		return "Player not have enough money!", NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", NO_ERROR
}
