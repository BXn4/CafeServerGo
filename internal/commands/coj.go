package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	FULL = "85"
)

func CoopJoin(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	if c.Player.GetActiveCoopID() != 0 {
		c.SendExtensionResponse("coj", "-1", "0", MAX_DONE)
		return nil
	}

	coopIDToJoin, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	coop, err := c.DB.GetCoop(coopIDToJoin)
	if err != nil {
		return fmt.Errorf("Cannot get and join coop!")
	}

	coopInfo, err := utils.GetCoop(coop.ActiveCoop)
	if err != nil {
		return err
	}

	if c.Player.GetLevel() > coopInfo.MaxLevel {
		c.SendExtensionResponse("coj", "-1", HIGH_LEVEL, strconv.Itoa(coopInfo.MaxLevel))
		return nil
	}

	if slices.Contains(coop.Members, c.Player.ID) {
		c.SendExtensionResponse("coj", "-1", "0", MAX_DONE)
		return nil
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		dishID, err := strconv.Atoi(dishRequirement[0])

		dishInfo, err := utils.GetDish(dishID)

		if err != nil {
			return nil
		}

		if dishInfo.Level > c.Player.GetLevel() {
			c.SendExtensionResponse("coj", "-1", LOW_LEVEL)
			return nil
		}
	}

	if len(coop.Members) >= coopInfo.MaxMembers {
		c.SendExtensionResponse("coj", "-1", FULL)
		return nil
	}

	if !coop.GetIsActive() {
		// if its ended
		return nil
	}

	coop.Join(c.Player.ID)

	var playersList []string

	for _, playerID := range coop.Members {
		player, err := c.DB.GetPlayer(playerID)

		if err != nil {
			playersList = append(playersList, "")
		} else {
			playersList = append(playersList, strconv.Itoa(player.ID)+"+"+strconv.Itoa(player.XP)+"+"+player.Avatar.String())
		}
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
