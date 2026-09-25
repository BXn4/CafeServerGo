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

// min level 4

func SendFriendsAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	friends := c.Player.GetFriends()
	if len(friends) == 0 {
		c.SendExtensionResponse("bga", "-1", "0", "")
		return nil
	}

	var friendsStr []string
	for _, f := range friends {

		var p *player.Player
		// Check if online
		item, err := gm.GetClient(f)
		if err == nil {
			oc := item.(*client.Client)
			p = oc.Player
		} else {
			// Get it from db
			p, err = c.DB.GetPlayer(f)
			if err != nil {
				return fmt.Errorf("Player %v not in db: %v", f, err)
			}
		}
		avatar := p.GetAvatar()
		friendsStr = append(friendsStr, fmt.Sprintf("%v+%v+%v", c.Player.GetID(), p.GetXP(), avatar.String()))
	}
	c.SendExtensionResponse("bga", "-1", "0", strings.Join(friendsStr, "%"))
	return nil
}
