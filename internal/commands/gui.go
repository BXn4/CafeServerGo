package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// gui - USER INFO
func UserInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var err error
	name := req.Args[2]

	if c.Player == nil {
		c.Player, err = c.DB.GetPlayerByName(name)
		if err != nil {
			return err
		}
	}

	c.SendExtensionResponse(
		"gui",
		"-1",
		"0",
		strconv.Itoa(c.Player.ID),
		strconv.Itoa(c.Player.ID),
		"0",
		strconv.Itoa(c.Player.Cash),
		strconv.Itoa(c.Player.Gold),
		strconv.Itoa(c.Player.XP),
		strconv.Itoa(c.Player.InstantCookings),
		strconv.Itoa(c.Player.OpenJobs),
		"0",
		utils.If(c.Player.PlayedWheel, "1", "0"),
		"0",
		utils.If(c.Player.AllowFriendRequests, "1", "0"),
		utils.If(c.Player.AllowEmails, "1", "0"),
		utils.If(c.Player.EmailVerified, "1", "0"),
		strconv.Itoa(c.Player.NewGifts),
		c.Player.Avatar.String(),
	)
	return nil
}
