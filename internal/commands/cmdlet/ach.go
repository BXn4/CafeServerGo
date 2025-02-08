package cmdlet

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"fmt"
	"strconv"
)

func SetAchivement(c *client.Client, gm *managers.GameManager, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("/ach <achivement id> <proggression>")
	}

	achivementID, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	proggression, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	c.Player.SetAchievement(achivementID, proggression)

	return nil
}
