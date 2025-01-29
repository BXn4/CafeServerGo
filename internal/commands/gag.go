package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func DailyGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Get last time gifts were refreshed
	refreshTime, err := c.DB.GetGiftRefreshTime(c.Player.ID)
	if err != nil {
		return err
	}

	// Sendable gifts
	var gifts string

	// If time is more tha 24 hours reset the time
	if time.Now().Sub(*refreshTime) < time.Hour*24 { // DEBUG: < (fliped)
		c.DB.ResetDailyLogin(c.Player.ID)

		// Get all possible gifts (TODO: add more items)
		possibleGifts, err := GetPossibleGifts()
		if err != nil {
			return err
		}

		// Generate new gifts
		giftsStr := []string{}
		for i := 0; i < 6; i++ {
			// Get gift and amount
			choice := rand.Intn(len(possibleGifts))
			gift := possibleGifts[choice]
			amount := rand.Intn(10) + 1

			// Add to gifts list
			giftStr := fmt.Sprintf("%v+%v", gift, amount)
			giftsStr = append(giftsStr, giftStr)
		}
		gifts = strings.Join(giftsStr, "#")

		// Save new Gifts
		c.DB.UpdateSendableGifts(c.Player.ID, gifts)
	} else {
		// Get gifts
		gifts, err = c.DB.GetSendableGifts(c.Player.ID)
		if err != nil {
			return err
		}
	}

	println("sendable gifts: ", gifts)

	c.SendExtensionResponse("gag", "-1", "0", gifts)
	return nil
}

func GetPossibleGifts() ([]string, error) {
	var possibleGifts []string

	// Add fancies to possible gifts
	fancies, err := utils.GetItems("fancy")
	if err != nil {
		return nil, err
	}

	for _, item := range fancies {
		possibleGifts = append(possibleGifts, strconv.Itoa(item.ID))
	}

	// Add decorations to possible gifts
	decos, err := utils.GetItems("deco")
	if err != nil {
		return nil, err
	}

	for _, item := range decos {
		possibleGifts = append(possibleGifts, strconv.Itoa(item.ID))
	}

	// Add dishes to possible gifts
	dishes, err := utils.GetItems("dish")
	if err != nil {
		return nil, err
	}

	for _, item := range dishes {
		possibleGifts = append(possibleGifts, strconv.Itoa(item.ID))
	}

	return possibleGifts, nil
}
