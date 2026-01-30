package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	commands.RegisterCommand(requests.C2S_COOP_LEAVE,
		commands.CommandConfig{
			Name:         "CoopLeave",
			Identifier:   responses.S2C_COOP_LEAVE,
			Description:  "Leave an coop",
			Args:         "{}",
			MinArgs:      2,
			MaxArgs:      2,
			FeatureLevel: 5,
		},
		CoopLeaveValidator,
		CoopLeave,
		CoopLeaveDBSaver,
	)
}

// col - CoopLeave
func CoopLeave(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	coop, _ := c.DB.GetCoop(c.Player.GetCoopID())
	coop.Leave(c.Player.GetID())

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

	c.Player.SetCoopID(0)
	c.SendExtensionResponse(cm.Identifier, "-1", "0")

	return nil
}

func CoopLeaveValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet reached coops!", commands.CONVERT_ERROR
	}

	if c.Player.GetCoopID() == 0 {
		return "Player is not in coop!", commands.NOT_DECLARED
	}

	_, err := c.DB.GetCoop(c.Player.GetCoopID())
	if err != nil {
		return "Cannot get coop!", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func CoopLeaveDBSaver(c *client.Client) error {
	c.DB.UpdateCoopID(c.Player.GetID(), c.Player.GetCoopID())

	return nil
}
