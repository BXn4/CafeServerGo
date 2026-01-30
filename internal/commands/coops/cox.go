package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"time"
)

func init() {
	commands.RegisterCommand(requests.C2S_COOP_EXTEND,
		commands.CommandConfig{
			Name:         "CoopExtend",
			Identifier:   responses.S2C_COOP_EXTEND,
			Description:  "Extend the coop leftover time",
			Args:         "{coop}",
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 5,
		},
		CoopExtendValidator,
		CoopExtend,
		CoopExtendDBSaver,
	)
}

func CoopExtend(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	/*coopID, err := strconv.Atoi(req.Args[1]) not need
	if err != nil {
		return err
		} */

	coop, _ := c.DB.GetCoop(c.Player.GetCoopID())
	c.Player.AddGold(-1)
	coop.AddExtend()
	coop.End = coop.End.Add(time.Duration(20) * time.Hour) // the time is fixed.
	// adds 20 hour to bronze, and we calculate the levels from the bronze time

	c.DB.SaveCoop(&coop)

	c.SendExtensionResponse(cm.Identifier, "-1", "0%", coop.GetCoop().AsResponse())

	return nil
}

func CoopExtendValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
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

	if c.Player.GetCoopID() == 0 {
		return "Cant extend coop, because the player not in coop!", commands.NOT_DECLARED
	}

	_, err := c.DB.GetCoop(c.Player.GetCoopID())
	if err != nil {
		return "Cant get coop detail from db!", commands.NOT_DECLARED
	}

	if c.Player.GetGold() < balancing.BalancingConstants.CoopExpansionGold {
		return "Player not have enough money!", commands.NOT_ENOUGHT_MONEY
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func CoopExtendDBSaver(c *client.Client) error {
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())

	return nil
}
