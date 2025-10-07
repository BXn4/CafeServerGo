package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
)

// ChangeAvatar - cha command
func ChangeAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	newAvatar := avatar.NewAvatarFromString(req.Args[2])

	if c.Player.AvatarChanged && c.Player.IsRegistered {
		if newAvatar.Apperance() == c.Player.Avatar.Apperance() {
			c.SendExtensionResponse("cha", "1", "-1")
			return fmt.Errorf("The new avatar is the same as the old!")
		}

		if c.Player.Gold < 1 {
			c.SendExtensionResponse("cha", "-1", "4")
			return fmt.Errorf("Player does not have enough gold to change the avatar")
		}
		c.Player.Gold -= 1
		c.DB.UpdateGold(c.Player.ID, c.Player.Gold)
	}

	// If the player at the register dialog, allow them to change their avatar free. Just need to reduce gold in the game.
	if c.Player.IsRegistered {
		c.Player.AvatarChanged = true
		c.DB.UpdateAvatarChanged(c.Player.ID, c.Player.AvatarChanged)
	}

	if !newAvatar.IsValid() {
		c.SendExtensionResponse("cha", "1", "-1")
		return fmt.Errorf("Invalid avatar structure!")
	}

	newAvatar.Name = c.Player.Username
	c.Player.Avatar = *newAvatar
	c.DB.UpdateAvatar(c.Player.ID, c.Player.Avatar.Apperance())

	log.Debugf("Player new avatar is: %s", newAvatar.String())

	if c.Location != nil {
		c.Location.Broadcast("cha", "1", "0", c.Player.Avatar.String(), strconv.Itoa(c.Player.ID))
	} else {
		c.SendExtensionResponse("cha", "1", "0", c.Player.Avatar.String(), strconv.Itoa(c.Player.ID))
	}
	return nil
}
