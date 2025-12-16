package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func init() {
	RegisterCommand(requests.C2S_COOP_START,
		CommandConfig{
			Name:       "CoopStart",
			Identifier: responses.S2C_COOP_START,
			MinArgs:    3,
			MaxArgs:    3,
		},
		CoopStartValidator,
		CoopStart,
	)
}

// cos - CoopStart
func CoopStart(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	coopIDToStart, _ := strconv.Atoi(req.Args[2])
	coopInfo, _ := utils.GetCoop(coopIDToStart)

	end := time.Now().UTC().Add(time.Duration(coopInfo.Duration) * time.Minute)

	coopID, err := c.DB.CreateCoop(coopInfo.ID, c.Player.ID, end)

	if err != nil {
		return fmt.Errorf("Cannot register the coop in the db: %w", err)
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return fmt.Errorf("Cannot get the registered coop: %w", err)
	}

	c.Player.SetActiveCoopID(coop.ID)

	c.DB.UpdateCoopID(c.Player.ID, c.Player.GetActiveCoopID())
	c.DB.UpdateStartedCoop(c.Player.ID, c.Player.GetStartedCoop())

	c.SendExtensionResponse("cos", "-1", "0", strconv.Itoa(coop.Host), coop.GetCoop().AsResponse())
	return nil
}

func CoopStartValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.GetStartedCoop() {
		return "Player started a coop today already!", ERROR_COOP_JOIN_MAX_DONE
	}

	if c.Player.GetActiveCoopID() != 0 {
		return "Player is in coop already!", ERROR_COOP_JOIN_MAX_DONE
	}

	coopInfo, err := utils.GetCoop(coopID)
	if err != nil {
		return "Cant get coop info!", NOT_DECLARED
	}

	if c.Player.GetLevel() > coopInfo.MaxLevel {
		return "Player level too high for the active coop!", ERROR_COOP_JOIN_LEVEL_HIGH
	}

	if event.GetEvent() < coopInfo.Events {
		return "Invalid coop ID, because theres no holiday", NOT_DECLARED
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")
	minLevel := 0

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		if len(dishRequirement) == 0 {
			continue
		}

		dishID, err := strconv.Atoi(dishRequirement[0])
		if err != nil {
			return "Cant convert string to int!", CONVERT_ERROR
		}

		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Cant get dish info!", CONVERT_ERROR
		}

		if dishInfo.Level < minLevel {
			minLevel = dishInfo.Level
		}
	}

	if minLevel == math.MaxInt {
		minLevel = 0
	}

	if c.Player.GetLevel() < minLevel {
		return "Player level too low for the active coop!", ERROR_COOP_JOIN_LOW_LEVEL
	}

	return "Command ran without any errors.", NO_ERROR
}
