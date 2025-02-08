package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
	"strings"
)

// cal
func SendAchivements(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	achivements := map[int]int{}

	// If online get player
	other, err := gm.GetClient(id)
	if err == nil {
		otherClient := other.(*client.Client)
		achivements = otherClient.Player.GetAchivements()

	} else {
		// else get it from db
		p, err := c.DB.GetPlayer(id)
		if err != nil {
			c.SendExtensionResponse("cal", "-1", "0")
			return fmt.Errorf("Can not get player's achivements: %v", id)
		}
		achivements = p.GetAchivements()
	}

	var achivementsStr []string
	for id, value := range achivements {
		if id > 2000 {
			id -= 2001
		}
		achivementsStr = append(achivementsStr, fmt.Sprintf("%v+%v", id, value))
	}

	c.SendExtensionResponse("cal", "-1", "0",
		req.Args[2], // Player who wants it
		req.Args[2], // Player who owns the achivements
		strings.Join(achivementsStr, "#"),
	)

	return nil
}
