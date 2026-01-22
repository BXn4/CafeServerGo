package gifts

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

// min level: 7

func init() {
	commands.RegisterCommand(requests.C2S_GIFT_REMOVE,
		commands.CommandConfig{
			Name:         "GiftRemove",
			Identifier:   responses.S2C_GIFT_REMOVE,
			Description:  "Remove a gift from the gift list",
			Args:         "{giftID}",
			MinArgs:      3,
			MaxArgs:      3,
			FeatureLevel: 7,
		},
		RemoveGiftValidator,
		RemoveGift,
		RemoveGiftDBSaver,
	)
}

func RemoveGift(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {

	slot, _ := strconv.Atoi(req.Args[2])
	gifts := c.Player.GetGifts()
	gifts.RemoveGift(slot)
	c.SendExtensionResponse(cm.Identifier, "-1", "0", req.Args[2])

	SendPlayerGifts(req, c, gm)
	return nil
}

func RemoveGiftValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
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

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	gifts := c.Player.GetGifts()

	if gifts[slot] != nil {
		return "Gift not found!", commands.NOT_DECLARED
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func RemoveGiftDBSaver(c *client.Client) error {
	gifts := c.Player.GetGifts()

	c.DB.UpdateGifts(c.Player.GetID(), gifts.String())

	return nil
}
