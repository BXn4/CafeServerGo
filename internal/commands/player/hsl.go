/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/leaderboard"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_HIGHSCORE_LIST,
		commands.CommandConfig{
			Name:         "HighScoreList",
			Identifier:   responses.S2C_HIGHSCORE_LIST,
			Description:  "Ingame leaderboard",
			Args:         "{offset} {leaderboardSTR}",
			MinArgs:      4,
			MaxArgs:      4,
			FeatureLevel: 2,
		},
		HighscoreListValidator,
		HighscoreList,
		nil,
	)
}

// hsl - Highscore list
func HighscoreList(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	search := req.Args[2]
	orderBy := utils.If(req.Args[3] == "0", leaderboard.OrderByXP, leaderboard.OrderByLuxury)
	offset := 0

	if search == "-1" {
		search = strconv.Itoa(leaderboard.GetPlayerRankByID(c.Player.GetID(), orderBy))
	}

	offset = leaderboard.GetOffset(leaderboard.GetPlayerRank(search, orderBy))

	c.SendExtensionResponse(cm.Identifier, "-1", "0", strconv.Itoa(offset), leaderboard.GetLeaderBoard(search, orderBy))
	return nil
}

func HighscoreListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet unlocked the feature", commands.NOT_DECLARED
	}

	search := req.Args[2]
	if search != "-1" {
		_, err := strconv.Atoi(search)
		if err != nil {
			return "Cant convert string to int!", commands.CONVERT_ERROR
		}
	}
	orderBy, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}
	switch orderBy {
	case leaderboard.OrderByXP, leaderboard.OrderByLuxury:
	default:
		return "Invalid order parameter!", commands.INVALID_ARGS
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
