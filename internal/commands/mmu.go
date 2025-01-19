package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"math/rand"
	"strconv"
	"time"
)

// mmu - S2C_MINI_MUFFIN
func PlayMuffinGame(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	isWin := random.Float64() <= 1/3.0

	cash, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	gold, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	/* if player level not enough to play {
		c.SendExtensionResponse("mmu", "-1", "9")
		return nil
	} */

	// If the player modifies the package, dont allow to send us negative numbers
	if cash < 100 || gold < 0 {
		return nil
	}

	// Just dont enable this
	if c.Player.Cash < cash || c.Player.Gold < gold {
		return nil
	}

	c.Player.Cash -= cash
	c.Player.Cash -= gold

	switch isWin {
	case true:
		c.Player.Cash += cash * 2
		c.Player.Gold += gold * 2

		c.SendExtensionResponse("mmu", "-1", "0", "1", strconv.Itoa(cash), strconv.Itoa(gold))
	case false:
		c.Player.Cash -= cash
		c.Player.Gold -= gold

		c.SendExtensionResponse("mmu", "-1", "0", "0", strconv.Itoa(-cash), strconv.Itoa(-gold))
	}

	return nil
}
