package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func init() {
	RegisterCommand(requests.C2S_COOP_JOIN,
		CommandConfig{
			Name:       "CoopJoin",
			Identifier: responses.S2C_COOP_JOIN,
			MinArgs:    3,
			MaxArgs:    3,
		},
		CoopJoinValidator,
		CoopJoin,
	)
}

func CoopJoin(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	coopIDToJoin, _ := strconv.Atoi(req.Args[2])
	coop, _ := c.DB.GetCoop(coopIDToJoin)

	coop.Join(c.Player.ID)

	// TODO: MAKE IT BETTER. MAKE THE COOPS SAVED IN THE MEMORY!
	var playersList []string
	for _, playerID := range coop.Members {
		player, _ := c.DB.GetPlayer(playerID)
		playersList = append(playersList, strconv.Itoa(player.ID)+"+"+strconv.Itoa(player.XP)+"+"+player.Avatar.String())
	}

	coop.SetPlayersString(strings.Join(playersList, "%"))

	c.Player.SetActiveCoopID(coop.ID)

	c.DB.UpdateCoopID(c.Player.ID, c.Player.GetActiveCoopID())

	c.DB.SaveCoop(&coop)

	c.SendExtensionResponse("coj", "-1", "0", strconv.Itoa(c.Player.ID), coop.GetCoop().AsResponse())

	for memberID := range coop.Members {
		if memberID != c.Player.ID {
			item, err := gm.GetClient(memberID)
			var toClient *client.Client
			if err == nil {
				toClient = item.(*client.Client)

				toClient.SendExtensionResponse("coj", "-1", "0", strconv.Itoa(c.Player.ID), coop.GetCoop().AsResponse())
			}
		}
	}

	return nil
}

func CoopJoinValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.GetLevel() < 5 {
		return "Player not yet reached coops!", CONVERT_ERROR
	}

	if c.Player.GetActiveCoopID() != 0 {
		return "Player is in coop already!", ERROR_COOP_JOIN_MAX_DONE
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return "Cant get coop detail from db!", NOT_DECLARED
	}

	coopInfo, err := utils.GetCoop(coop.ActiveCoop)
	if err != nil {
		return "Cant get coop info!", NOT_DECLARED
	}

	if c.Player.GetLevel() > coopInfo.MaxLevel {
		return "Player level too high for the active coop!", ERROR_COOP_JOIN_LEVEL_HIGH
	}

	if slices.Contains(coop.Members, c.Player.ID) {
		return "Player is in coop already!", ERROR_COOP_JOIN_MAX_DONE
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		dishID, err := strconv.Atoi(dishRequirement[0])

		if err != nil {
			return "Cant convert string to int!", CONVERT_ERROR
		}

		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Cant get dishinfo!", CONVERT_ERROR
		}
		if dishInfo.Level > c.Player.GetLevel() {
			return "Player level too low for the active coop!", ERROR_COOP_JOIN_LOW_LEVEL
		}
	}

	if len(coop.Members) >= coopInfo.MaxMembers {
		return "Coop is full!", ERROR_COOP_JOIN_FULL
	}

	if !coop.GetIsActive() {
		return "Coop ended!", NOT_DECLARED
	}

	for _, playerID := range coop.Members {
		_, err := c.DB.GetPlayer(playerID)
		if err != nil {
			return "Cant get other coop members!", NOT_DECLARED
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
