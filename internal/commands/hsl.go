package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/leaderboard"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_HIGHSCORE_LIST,
		CommandConfig{
			Name:       "HighScoreList",
			Identifier: responses.S2C_HIGHSCORE_LIST,
			MinArgs:    4,
			MaxArgs:    4,
		},
		HighscoreListValidator,
		HighscoreList,
	)
}

// hsl - Highscore list
func HighscoreList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	search := req.Args[2]
	orderBy := utils.If(req.Args[3] == "0", leaderboard.OrderByXP, leaderboard.OrderByLuxury)
	offset := 0

	if search == "-1" {
		search = strconv.Itoa(leaderboard.GetPlayerRankByID(c.Player.ID, orderBy))
	}

	offset = leaderboard.GetOffset(leaderboard.GetPlayerRank(search, orderBy))

	c.SendExtensionResponse("hsl", "-1", "0", strconv.Itoa(offset), leaderboard.GetLeaderBoard(search, orderBy))
	return nil
}

func HighscoreListValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if c.Player.GetLevel() < 2 {
		return "Player not yet unlocked the feature", NOT_DECLARED
	}

	search := req.Args[2]
	if search != "-1" {
		_, err := strconv.Atoi(search)
		if err != nil {
			return "Cant convert string to int!", CONVERT_ERROR
		}
	}
	orderBy, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}
	switch orderBy {
	case leaderboard.OrderByXP, leaderboard.OrderByLuxury:
	default:
		return "Invalid order parameter!", INVALID_ARGS
	}

	return "Command ran without any errors.", NO_ERROR
}
