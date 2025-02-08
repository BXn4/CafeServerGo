package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

const (
	LOW_LEVEL  = "3"
	HIGH_LEVEL = "7"
	MAX_DONE   = "87"
)

// cos - CoopStart
func CoopStart(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Parse coop id
	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	// Get coop
	coop, err := utils.GetCoop(coopID)
	if err != nil {
		return err
	}

	// Check if player is in a coop
	if c.Player.CoopID > 0 {
		c.SendExtensionResponse("cos", "-1", MAX_DONE)
		return nil
	}

	// Check if player level is enough
	if coop.MaxLevel < c.Player.GetLevel() {
		c.SendExtensionResponse("cos", "-1", HIGH_LEVEL, strconv.Itoa(coop.MaxLevel))
		return nil
	}

	c.SendExtensionResponse("cos", "1", "0")
	return nil
}
