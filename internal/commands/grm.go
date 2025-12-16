package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
)

// min level: 7

func init() {
	RegisterCommand(requests.C2S_GIFT_REMOVE,
		CommandConfig{
			Name:       "GiftRemove",
			Identifier: responses.S2C_GIFT_REMOVE,
			MinArgs:    3,
			MaxArgs:    3,
		},
		RemoveGiftValidator,
		RemoveGift,
	)
}

func RemoveGift(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	slot, _ := strconv.Atoi(req.Args[2])
	c.Player.Gifts.RemoveGift(slot)
	c.SendExtensionResponse("gmg", "-1", "0", req.Args[2])

	c.DB.UpdateGifts(c.Player.ID, c.Player.Gifts.String())

	SendPlayerGifts(req, c, gm)
	return nil
}

func RemoveGiftValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	slot, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.Gifts[slot] != nil {
		return "Gift not found!", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
