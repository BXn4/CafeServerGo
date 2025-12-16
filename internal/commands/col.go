package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	RegisterCommand(requests.C2S_COOP_LEAVE,
		CommandConfig{
			Name:       "CoopLeave",
			Identifier: responses.S2C_COOP_LEAVE,
			MinArgs:    3,
			MaxArgs:    3,
		},
		CoopLeaveValidator,
		CoopLeave,
	)
}

// col - CoopLeave
func CoopLeave(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	coop, _ := c.DB.GetCoop(c.Player.GetActiveCoopID())
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

func CoopLeaveValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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
		return "Player is not in coop!", NOT_DECLARED
	}

	_, err := c.DB.GetCoop(c.Player.GetActiveCoopID())
	if err != nil {
		return "Cannot get coop!", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
