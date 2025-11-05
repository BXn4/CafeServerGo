package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/leaderboard"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// hsl - Highscore list
func SendHighscoreList(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	search := (req.Args[2])
	orderBy := utils.If(req.Args[3] == "0", leaderboard.OrderByXP, leaderboard.OrderByLuxury)
	offset := 0

	if search == "-1" {
		search = strconv.Itoa(leaderboard.GetPlayerRankByID(c.Player.ID, orderBy))
	}

	offset = leaderboard.GetOffset(leaderboard.GetPlayerRank(search, orderBy))

	c.SendExtensionResponse("hsl", "-1", "0", strconv.Itoa(offset), leaderboard.GetLeaderBoard(search, orderBy))
	return nil
}
