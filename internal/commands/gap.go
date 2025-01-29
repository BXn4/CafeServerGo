package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"slices"
	"strconv"
)

func GiftAllReadySendPlayers(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Get friends
	friends := c.Player.Friends

	// Get friends who got gifts
	friendsWithGifts, err := c.DB.GetFriendsWithGifts(c.Player.ID)
	if err != nil {
		return err
	}

	// Gather friends who didnt get gifts
	var friendsWithoutGifts []string
	for _, friend := range friends {
		idStr := strconv.Itoa(friend)

		// If got gift continue
		if slices.Contains(friendsWithGifts, idStr) {
			continue
		}

		friendsWithoutGifts = append(friendsWithoutGifts, strconv.Itoa(friend))
	}

	// DEBUG/TEST
	err = SendSocialFriendsAvatar(req, c, gm)
	if err != nil {
		return err
	}

	// Send friends who did not get gifts
	c.SendExtensionResponse("gap", "-1", "0",
		// 0 - Majd holnap küldhetsz ajándékot a barátaidnak.
		// 1 - Egy nap egy ajándékot adhatsz egy barátodat.
		"1", // canSendGifts
		// strings.Join(friendsWithoutGifts, "+"),
		fmt.Sprintf("%v+%v+%v", c.Player.ID, c.Player.GetXP(), c.Player.Avatar.String(c.Player.Username)),
	)
	return nil
}
