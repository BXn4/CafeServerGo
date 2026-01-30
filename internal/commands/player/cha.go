package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"

	"github.com/charmbracelet/log"
)

func init() {
	commands.RegisterCommand(requests.C2S_CHANGE_AVATAR,
		commands.CommandConfig{
			Name:        "ChangeAvatar",
			Identifier:  responses.S2C_CHANGE_AVATAR,
			Description: "Changing player avatar",
			Args:        "{playerAvatar} {playerID}",
			MinArgs:     3,
			MaxArgs:     3,
		},
		ChangeAvatarValidator,
		ChangeAvatar,
		ChangeAvatarDBSaver,
	)
}

// ChangeAvatar - cha command
func ChangeAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	newAvatar := avatar.NewAvatarFromString(req.Args[2])

	if c.Player.GetAvatarChanged() && c.Player.GetIsRegistered() {
		c.Player.AddGold(-1)
	}

	// If the player at the register dialog, allow them to change their avatar free. Just need to reduce gold in the game.
	if c.Player.GetIsRegistered() {
		c.Player.SetAvatarChanged(true)
	}

	newAvatar.Name = c.Player.GetUsername()
	c.Player.SetAvatar(*newAvatar)

	log.Debugf("Player new avatar is: %s", newAvatar.String())

	avatar := c.Player.GetAvatar()

	if c.Location != nil {
		c.Location.Broadcast(cm.Identifier, "1", "0", avatar.String(), strconv.Itoa(c.Player.GetID()))
	} else {
		c.SendExtensionResponse(cm.Identifier, "1", "0", avatar.String(), strconv.Itoa(c.Player.GetID()))
	}
	return nil
}

func ChangeAvatarValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	newAvatar := avatar.NewAvatarFromString(req.Args[2])

	if newAvatar == nil {
		return "Cant parse avatar from the string, avatar is invalid!", commands.NOT_DECLARED
	}

	if !newAvatar.IsValid() {
		return "Cant parse avatar from the string, avatar is invalid!", commands.NOT_DECLARED
	}

	if c.Player.GetAvatarChanged() && c.Player.GetIsRegistered() {
		avatar := c.Player.GetAvatar()
		if newAvatar.Apperance() == avatar.Apperance() {
			return "The avatar is same as the old avatar!", commands.NOT_DECLARED
		}

		if c.Player.GetGold() < 1 {
			return "Player not have enough cash to update avatar", commands.NOT_ENOUGHT_MONEY
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func ChangeAvatarDBSaver(c *client.Client) error {
	avatar := c.Player.GetAvatar()

	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateAvatarChanged(c.Player.GetID(), c.Player.GetAvatarChanged())
	c.DB.UpdateAvatar(c.Player.GetID(), avatar.String())

	return nil
}
