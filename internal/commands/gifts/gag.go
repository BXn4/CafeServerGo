package gifts

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/gift"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math/rand"
	"time"
)

// min level: 7

func init() {
	commands.RegisterCommand(requests.C2S_GIFT_SENDABLEGIFTS,
		commands.CommandConfig{
			Name:         "DailyGifts",
			Description:  "Not fully known yet",
			Args:         "{gifts}",
			Identifier:   responses.S2C_GIFT_SENDABLEGIFTS,
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 7,
		},
		DailyGiftsValidator,
		DailyGifts,
		nil,
	)
}

func DailyGifts(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {

	// Sendable gifts
	gifts := gift.GiftList{}

	// If time is more tha 24 hours reset the time
	if time.Now().Sub(c.Player.GetGiftRefreshTime()) < time.Hour*24 { // DEBUG: < (fliped)
		// Reset daily login
		c.Player.SetDailyLogin(time.Now())

		// Get all possible gifts (TODO: add more items)
		possibleGifts, _ := GetPossibleGifts()

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
		c.Player.SetSendableGifts(gifts)
	} else {
		// Get gifts
		gifts = c.Player.GetSendableGifts()
	}

	c.SendExtensionResponse(cm.Identifier, "-1", "0", gifts.String())
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

func DailyGiftsValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
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

	_, err := GetPossibleGifts()
	if err != nil {
		return "Cant get possible gifts!", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}
