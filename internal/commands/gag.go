package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/gift"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"math/rand"
	"time"
)

func DailyGifts(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Sendable gifts
	gifts := gift.GiftList{}

	// If time is more tha 24 hours reset the time
	if time.Now().Sub(c.Player.GiftRefreshTime) < time.Hour*24 { // DEBUG: < (fliped)
		// Reset daily login
		c.Player.DailyLogin = time.Now()

		// Get all possible gifts (TODO: add more items)
		possibleGifts, err := GetPossibleGifts()
		if err != nil {
			return err
		}

		// Generate new gifts
		for i := 0; i < 6; i++ {
			// Get gift and amount
			choice := rand.Intn(len(possibleGifts))
			id := possibleGifts[choice]
			amount := rand.Intn(10) + 1

			g := &gift.Gift{
				ID:     id,
				Amount: amount,
			}

			gifts = append(gifts, g)
		}

		// Save new Gifts (update sendable gifts)
		c.Player.SendableGifts = gifts
	} else {
		// Get gifts
		gifts = c.Player.SendableGifts
	}

	c.SendExtensionResponse("gag", "-1", "0", gifts.String())
	return nil
}

func GetPossibleGifts() ([]int, error) {
	var possibleGifts []int

	// Add fancies to possible gifts
	fancies, err := utils.GetItems("fancy")
	if err != nil {
		return nil, err
	}

	for _, item := range fancies {
		possibleGifts = append(possibleGifts, item.ID)
	}

	// Add decorations to possible gifts
	decos, err := utils.GetItems("deco")
	if err != nil {
		return nil, err
	}

	for _, item := range decos {
		possibleGifts = append(possibleGifts, item.ID)
	}

	// Add dishes to possible gifts
	dishes, err := utils.GetItems("dish")
	if err != nil {
		return nil, err
	}

	for _, item := range dishes {
		possibleGifts = append(possibleGifts, item.ID)
	}

	return possibleGifts, nil
}
