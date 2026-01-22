package player

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

	avatar := c.Player.GetAvatar()

	c.SendExtensionResponse(
		"gui",
		"-1",
		"0",
		strconv.Itoa(c.Player.GetID()),
		strconv.Itoa(c.Player.GetID()),
		"0",
		strconv.Itoa(c.Player.GetCash()),
		strconv.Itoa(c.Player.GetGold()),
		strconv.Itoa(c.Player.GetXP()),
		strconv.Itoa(c.Player.GetInstantCookingsUsed()),
		strconv.Itoa(c.Player.GetOpenJobs()),
		"0",
		utils.If(c.Player.GetPlayedWheel(), "1", "0"),
		utils.If(c.Player.GetAvatarChanged(), "1", "0"),
		utils.If(c.Player.GetAllowFriendRequests(), "1", "0"),
		utils.If(c.Player.GetAllowEmails(), "1", "0"),
		utils.If(c.Player.GetEmailVerified(), "1", "0"),
		strconv.Itoa(len(c.Player.GetGifts())),
		avatar.String(),
	)
	return nil
}
