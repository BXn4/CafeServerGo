/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package friends

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"fmt"
	"strings"
)

/*
This file is temporary,
this was created just

*/

func SendSocialFriendsAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	friends := c.Player.GetFriends()
	if len(friends) == 0 {
		c.SendExtensionResponse("sbs", "-1", "0", "")
		return nil
	}

	var friendsStr []string
	for _, f := range friends {

		var player *player.Player
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
		avatar := player.GetAvatar()
		friendsStr = append(friendsStr, fmt.Sprintf("%v|%v|%v|%v", pln, c.Player.GetID(), player.GetXP(), avatar.String()))
	}
	c.SendExtensionResponse("sbs", "-1", "0", strings.Join(friendsStr, "||"))
	return nil
}
