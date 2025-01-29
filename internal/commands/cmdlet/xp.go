package cmdlet

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"fmt"
	"strconv"
)

func SetXP(c *client.Client, gm *managers.GameManager, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("/xp <player name> <xp amount>")
	}

	name := args[0]
	xpAmount, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	client, err := gm.GetClientByName(name)
	if err != nil {
		return err
	}
	client.Player.SetXP(xpAmount)

	return nil
}
