package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	_ "cafego/internal/objects"
	"cafego/internal/types/requests"
	"strconv"
)

// gui - USER INFO
func UserInfo(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
	var err error
	name := req.Args[2]

	if c.Player == nil {
		c.Player, err = c.DB.GetPlayerByName(name)
		if err != nil {
			return err
		}
	}

	playedWheel := "0"
	if c.Player.PlayedWheel {
		playedWheel = "1"
	}

	allowFriendRequests := "0"
	if c.Player.AllowFriendRequests {
		allowFriendRequests = "1"
	}

	allowEmails := "0"
	if c.Player.AllowEmails {
		allowEmails = "1"
	}

	emailVerified := "0"
	if c.Player.EmailVerified {
		emailVerified = "1"
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
		playedWheel,
		"0",
		allowFriendRequests,
		allowEmails,
		emailVerified,
		strconv.Itoa(c.Player.NewGifts),
		c.Player.Avatar.String(),
	)
	return nil
}
