package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"slices"
)

// min level: 7

func init() {
	RegisterCommand(requests.C2S_GIFT_SENDABLEGIFTS,
		CommandConfig{
			Name:       "SendableGifts",
			Identifier: responses.S2C_GIFT_SENDABLEGIFTS,
			MinArgs:    0,
			MaxArgs:    0,
		},
		SendableGiftsValidator,
		SendableGifts,
	)
}

func SendableGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Get friends
	friends := c.Player.Friends
	friendsWithGifts := []int(c.Player.FriendsWithGifts)

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
	err := SendSocialFriendsAvatar(req, c, gm)
	if err != nil {
		return err
	}

	// Send friends who did not get gifts
	c.SendExtensionResponse("gap", "-1", "0",
		// 0 - Majd holnap küldhetsz ajándékot a barátaidnak.
		// 1 - Egy nap egy ajándékot adhatsz egy barátodat.
		"1", // canSendGifts
		fmt.Sprintf("%v+%v+%v", c.Player.ID, c.Player.GetXP(), c.Player.Avatar.String()),
	)
	return nil
}

func SendableGiftsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if true {
		return "NOT YET FINISHED!", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
