/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package gifts

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"slices"
)

// min level: 7

func init() {
	commands.RegisterCommand(requests.C2S_GIFT_SENDABLEGIFTS,
		commands.CommandConfig{
			Name:         "SendableGifts",
			Description:  "Sending gifts. Only for social login!",
			Args:         "{playerID} {playerXP} {playerAvatar}",
			Identifier:   responses.S2C_GIFT_SENDABLEGIFTS,
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 7,
		},
		SendableGiftsValidator,
		SendableGifts,
		nil,
	)
}

func SendableGifts(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	// Get friends
	friends := c.Player.GetFriends()
	friendsWithGifts := []int(c.Player.GetFriendsWithGifts())

	// Gather friends who didnt get gifts
	var friendsWithoutGifts []int
	for _, friend := range friends {
		// If got gift continue
		if slices.Contains(friendsWithGifts, friend) {
			continue
		}

		friendsWithoutGifts = append(friendsWithoutGifts, friend)
	}

	// DEBUG/TEST
	/* err := SendSocialFriendsAvatar(req, c, gm)
	if err != nil {
		return err
	} */

	avatar := c.Player.GetAvatar()

	// Send friends who did not get gifts
	c.SendExtensionResponse(cm.Identifier, "-1", "0",
		// 0 - Majd holnap küldhetsz ajándékot a barátaidnak.
		// 1 - Egy nap egy ajándékot adhatsz egy barátodat.
		"1", // canSendGifts
		fmt.Sprintf("%v+%v+%v", c.Player.GetID(), c.Player.GetXP(), avatar.String()),
	)
	return nil
}

func SendableGiftsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet unlocked the feature!", commands.NOT_DECLARED
	}

	if true {
		return "NOT YET FINISHED!", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
