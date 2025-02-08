package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"fmt"
	"strings"
)

func SendFriendsAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	friends := c.Player.Friends
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
		friendsStr = append(friendsStr, fmt.Sprintf("%v+%v+%v", c.Player.ID, p.GetXP(), p.Avatar.String()))
	}
	c.SendExtensionResponse("bga", "-1", "0", strings.Join(friendsStr, "%"))
	return nil
}
