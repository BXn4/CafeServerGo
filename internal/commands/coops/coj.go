package coops

import (
	"cafego/internal/client"
	"cafego/internal/commands"
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
	commands.RegisterCommand(requests.C2S_COOP_JOIN,
		commands.CommandConfig{
			Name:         "CoopJoin",
			Identifier:   responses.S2C_COOP_JOIN,
			Description:  "Joining to a coop",
			Args:         "{playerID} {coop}",
			MinArgs:      3,
			MaxArgs:      3,
			FeatureLevel: 5,
		},
		CoopJoinValidator,
		CoopJoin,
		CoopJoinDBSaver,
	)
}

func CoopJoin(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	coopIDToJoin, _ := strconv.Atoi(req.Args[2])
	coop, _ := c.DB.GetCoop(coopIDToJoin)

	coop.Join(c.Player.GetID())

	// TODO: MAKE IT BETTER. MAKE THE COOPS SAVED IN THE MEMORY!
	var playersList []string
	for _, playerID := range coop.Members {
		player, _ := c.DB.GetPlayer(playerID)
		avatar := player.GetAvatar()
		playersList = append(playersList, strconv.Itoa(player.GetID())+"+"+strconv.Itoa(player.GetXP())+"+"+avatar.String())
	}

	coop.SetPlayersString(strings.Join(playersList, "%"))

	c.Player.SetCoopID(coop.ID)

	c.DB.SaveCoop(&coop)

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(c.Player.GetID()), coop.GetCoop().AsResponse())

	for memberID := range coop.Members {
		if memberID != c.Player.GetID() {
			item, err := gm.GetClient(memberID)
			var toClient *client.Client
			if err == nil {
				toClient = item.(*client.Client)

				toClient.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(c.Player.GetID()), coop.GetCoop().AsResponse())
			}
		}
	}

	return nil
}

func CoopJoinValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	coopID, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet reached coops!", commands.NOT_DECLARED
	}

	if c.Player.GetCoopID() != 0 {
		return "Player is in coop already!", commands.ERROR_COOP_JOIN_MAX_DONE
	}

	coop, err := c.DB.GetCoop(coopID)
	if err != nil {
		return "Cant get coop detail from db!", commands.NOT_DECLARED
	}

	coopInfo, err := utils.GetCoop(coop.ActiveCoop)
	if err != nil {
		return "Cant get coop info!", commands.NOT_DECLARED
	}

	if c.Player.GetLevel() > coopInfo.MaxLevel {
		return "Player level too high for the active coop!", commands.ERROR_COOP_JOIN_LEVEL_HIGH
	}

	if slices.Contains(coop.Members, c.Player.GetID()) {
		return "Player is in coop already!", commands.ERROR_COOP_JOIN_MAX_DONE
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		dishID, err := strconv.Atoi(dishRequirement[0])

		if err != nil {
			return "Cant convert string to int!", commands.CONVERT_ERROR
		}

		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Cant get dishinfo!", commands.CONVERT_ERROR
		}
		if dishInfo.Level > c.Player.GetLevel() {
			return "Player level too low for the active coop!", commands.ERROR_COOP_JOIN_LOW_LEVEL
		}
	}

	if len(coop.Members) >= coopInfo.MaxMembers {
		return "Coop is full!", commands.ERROR_COOP_JOIN_FULL
	}

	if !coop.GetIsActive() {
		return "Coop ended!", commands.NOT_DECLARED
	}

	for _, playerID := range coop.Members {
		_, err := c.DB.GetPlayer(playerID)
		if err != nil {
			return "Cant get other coop members!", commands.NOT_DECLARED
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func CoopJoinDBSaver(c *client.Client) error {
	c.DB.UpdateCoopID(c.Player.GetID(), c.Player.GetCoopID())

	return nil
}
