package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"fmt"
	"strings"
)

/*
This file is temporary,
this was created just

*/

func SendSocialFriendsAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	friends := c.Player.Friends
	if len(friends) == 0 {
		c.SendExtensionResponse("sbs", "-1", "0", "")
		return nil
	}

	var friendsStr []string
	for _, f := range friends {

		var player *objects.Player
		// Check if online
		item, err := gm.GetClient(f)
		if err == nil {
			oc := item.(*client.Client)
			player = oc.Player
		} else {
			// Get it from db
			player, err = c.DB.GetPlayer(f)
			if err != nil {
				return fmt.Errorf("Player %v not in db: %v", f, err)
			}
		}
		pln := "2"
		friendsStr = append(friendsStr, fmt.Sprintf("%v|%v|%v|%v", pln, c.Player.ID, player.GetXP(), player.Avatar.String(player.Username)))
	}
	c.SendExtensionResponse("sbs", "-1", "0", strings.Join(friendsStr, "||"))
	return nil
}
