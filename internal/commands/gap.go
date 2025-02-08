package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"slices"
)

func GiftAllReadySendPlayers(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

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
