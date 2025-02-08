package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

// pin
func SendAchivements(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	achivements := ""

	// If online get player
	other, err := gm.GetClient(id)
	if err == nil {
		otherClient := other.(*client.Client)
		achivements = otherClient.Player.BuildAchievement()
	} else {
		// else get it from db
		p, err := c.DB.GetPlayer(id)
		if err != nil {
			c.SendExtensionResponse("cal", "-1", "0")
			return fmt.Errorf("Can not get player's achivements: %v", id)
		}
		achivements = p.BuildAchievement()
	}

	if achivements == "" {
		c.SendExtensionResponse("cal", "-1", "0")
		return fmt.Errorf("Can not get player's achivements: %v", id)
	}

	c.SendExtensionResponse("cal", "-1", "0",
		req.Args[2], // Player who wants it
		req.Args[2], // Player who owns the achivements
		achivements,
	)

	return nil
}
