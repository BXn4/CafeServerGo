package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	LOW_LEVEL  = "3"
	HIGH_LEVEL = "7"
	MAX_DONE   = "87"
)

// cos - CoopStart
func CoopStart(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Parse coop id
	coopIDToStart, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	// Get coop
	coopInfo, err := utils.GetCoop(coopIDToStart)
	if err != nil {
		return err
	}

	// Check if player was is in a coop today
	if c.Player.GetStartedCoop() || c.Player.GetActiveCoopID() != 0 {
		c.SendExtensionResponse("cos", "-1", MAX_DONE)
		return nil
	}

	// Check if player level is enough
	if c.Player.GetLevel() > coopInfo.MaxLevel {
		c.SendExtensionResponse("cos", "-1", HIGH_LEVEL, strconv.Itoa(coopInfo.MaxLevel))
		return nil
	}

	if event.GetEvent() < coopInfo.Events {
		return fmt.Errorf("Invalid coop ID:, because theres no holiday")
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	minLevel := math.MaxInt

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		if len(dishRequirement) == 0 {
			continue
		}

		dishID, err := strconv.Atoi(dishRequirement[0])
		if err != nil {
			return nil
		}

		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return nil
		}

		if dishInfo.Level < minLevel {
			minLevel = dishInfo.Level
		}
	}

	if minLevel == math.MaxInt {
		minLevel = 0
	}

	if c.Player.GetLevel() < minLevel {
		c.SendExtensionResponse("cos", "-1", LOW_LEVEL)
		return nil
	}

	end := time.Now().UTC().Add(time.Duration(coopInfo.Duration) * time.Minute)

	coopID, err := c.DB.CreateCoop(coopInfo.ID, c.Player.ID, end)

	if err != nil {
		return fmt.Errorf("Cannot create a coop: %d", err)
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return fmt.Errorf("Cannot get a coop: %d", err)
	}

	c.Player.SetActiveCoopID(coop.ID)

	c.DB.UpdateCoopID(c.Player.ID, c.Player.GetActiveCoopID())
	c.DB.UpdateStartedCoop(c.Player.ID, c.Player.GetStartedCoop())

	c.SendExtensionResponse("cos", "-1", "0", strconv.Itoa(coop.Host), coop.GetCoop().AsResponse())
	return nil
}
