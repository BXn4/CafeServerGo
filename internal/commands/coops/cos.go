package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
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
	commands.RegisterCommand(requests.C2S_COOP_START,
		commands.CommandConfig{
			Name:         "CoopStart",
			Identifier:   responses.S2C_COOP_START,
			Description:  "Start a coop",
			Args:         "{hostID} {coop}",
			MinArgs:      3,
			MaxArgs:      3,
			FeatureLevel: 5,
		},
		CoopStartValidator,
		CoopStart,
		CoopDBSaver,
	)
}

// cos - CoopStart
func CoopStart(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	coopIDToStart, _ := strconv.Atoi(req.Args[2])
	coopInfo, _ := utils.GetCoop(coopIDToStart)

	end := time.Now().UTC().Add(time.Duration(coopInfo.Duration) * time.Minute)

	coopID, err := c.DB.CreateCoop(coopInfo.ID, c.Player.GetID(), end)

	if err != nil {
		return fmt.Errorf("Cannot register the coop in the db: %w", err)
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return fmt.Errorf("Cannot get the registered coop: %w", err)
	}

	c.Player.SetCoopID(coop.ID)

	c.SendExtensionResponse("cos", "-1", "0", strconv.Itoa(coop.Host), coop.GetCoop().AsResponse())
	return nil
}

func CoopStartValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < 5 {
		return "Player not yet reached coops!", commands.CONVERT_ERROR
	}

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if c.Player.GetIsStartedCoop() {
		return "Player started a coop today already!", commands.ERROR_COOP_JOIN_MAX_DONE
	}

	if c.Player.GetCoopID() != 0 {
		return "Player is in coop already!", commands.ERROR_COOP_JOIN_MAX_DONE
	}

	coopInfo, err := utils.GetCoop(coopID)
	if err != nil {
		return "Cant get coop info!", commands.NOT_DECLARED
	}

	if c.Player.GetLevel() > coopInfo.MaxLevel {
		return "Player level too high for the active coop!", commands.ERROR_COOP_JOIN_LEVEL_HIGH
	}

	if event.GetEvent() < coopInfo.Events {
		return "Invalid coop ID, because theres no holiday", commands.NOT_DECLARED
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
			return "Cant convert string to int!", commands.CONVERT_ERROR
		}

		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Cant get dish info!", commands.CONVERT_ERROR
		}

		if dishInfo.Level < minLevel {
			minLevel = dishInfo.Level
		}
	}

	if minLevel == math.MaxInt {
		minLevel = 0
	}

	if c.Player.GetLevel() < minLevel {
		return "Player level too low for the active coop!", commands.ERROR_COOP_JOIN_LOW_LEVEL
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func CoopDBSaver(c *client.Client) error {
	c.DB.UpdateCoopID(c.Player.GetID(), c.Player.GetCoopID())
	c.DB.UpdateStartedCoop(c.Player.GetID(), c.Player.GetIsStartedCoop())

	return nil
}
