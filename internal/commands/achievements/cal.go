/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package achievements

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_ACHIEVEMENT_LIST,
		commands.CommandConfig{
			Name:        "AchievementList",
			Identifier:  responses.S2C_CAFE_ACHIEVEMENT_LIST,
			Description: "Sending the achivement list",
			Args:        "{Location from leaderboard or not 0/1} {FromPlayerID} {TargetPlayerID}#{AchievementID+AchivementProgress}#...",
			MinArgs:     4,
			MaxArgs:     4,
		},
		AchievementListValidator,
		AchievementList,
		nil,
	)
}

// cal
func AchievementList(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	id, _ := strconv.Atoi(req.Args[2])
	achivements := map[int]int{}

	// If online get player
	other, err := gm.GetClient(id)
	if err == nil {
		otherClient := other.(*client.Client)
		achivements = otherClient.Player.GetAchivements()

	} else {
		// else get it from db
		p, _ := c.DB.GetPlayer(id)
		achivements = p.GetAchivements()
	}

	var achivementsStr []string
	for id, value := range achivements {
		if id > 2000 {
			id -= 2001
		}
		achivementsStr = append(achivementsStr, fmt.Sprintf("%v+%v", id, value))
	}

	c.SendExtensionResponse(cm.Identifier, "-1", "0",
		req.Args[2], // Player who wants it
		req.Args[2], // Player who owns the achivements
		strings.Join(achivementsStr, "#"),
	)

	return nil
}

func AchievementListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	id, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	_, err = gm.GetClient(id)
	if err == nil {
	} else {
		// else get it from db
		_, err := c.DB.GetPlayer(id)
		if err != nil {
			return "Cant get the player achievements!", commands.NOT_DECLARED
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
