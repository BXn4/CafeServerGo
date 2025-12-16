package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
)

func init() {
	RegisterCommand(requests.C2S_CHANGE_AVATAR,
		CommandConfig{
			Name:       "ChangeAvatar",
			Identifier: responses.S2C_CHANGE_AVATAR,
			MinArgs:    3,
			MaxArgs:    3,
		},
		ChangeAvatarValidator,
		ChangeAvatar,
	)
}

// ChangeAvatar - cha command
func ChangeAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	newAvatar := avatar.NewAvatarFromString(req.Args[2])

	if c.Player.AvatarChanged && c.Player.IsRegistered {
		c.Player.AddGold(-1)
		c.DB.UpdateGold(c.Player.ID, c.Player.Gold)
	}

	// If the player at the register dialog, allow them to change their avatar free. Just need to reduce gold in the game.
	if c.Player.IsRegistered {
		c.Player.AvatarChanged = true
		c.DB.UpdateAvatarChanged(c.Player.ID, c.Player.AvatarChanged)
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

func ChangeAvatarValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	newAvatar := avatar.NewAvatarFromString(req.Args[2])

	if newAvatar == nil {
		return "Cant parse avatar from the string, avatar is invalid!", NOT_DECLARED
	}

	if !newAvatar.IsValid() {
		return "Cant parse avatar from the string, avatar is invalid!", NOT_DECLARED
	}

	if c.Player.AvatarChanged && c.Player.IsRegistered {
		if newAvatar.Apperance() == c.Player.Avatar.Apperance() {
			return "The avatar is same as the old avatar!", NOT_DECLARED
		}

		if c.Player.Gold < 1 {
			return "Player not have enough cash to update avatar", NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
