package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	RegisterCommand(requests.C2S_MINI_MUFFIN,
		CommandConfig{
			Name:       "MuffinMinigame",
			Identifier: responses.S2C_MINI_MUFFIN,
			MinArgs:    4,
			MaxArgs:    4,
		},
		MuffinGameValidator,
		MuffinGame,
	)
}

// 10 -> gold bet // min level
//
// MAX GOLD TO BET: 10000 - MAX CASH TO BET: 9999999 (from client)

// mmu - S2C_MINI_MUFFIN
func MuffinGame(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	isWin := random.Float64() <= 1/3.0

	cash, _ := strconv.Atoi(req.Args[2])
	gold, _ := strconv.Atoi(req.Args[3])

	c.Player.AddCash(-cash)
	c.Player.AddGold(-gold)

	switch isWin {
	case true:
		c.Player.AddCash(cash * 2)
		c.Player.AddGold(gold * 2)
		c.Player.UpdateAchivementMuffinmanCash(cash * 2)
		c.Player.UpdateAchivementMuffinmanGold(gold * 2)

		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

		c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
		c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

		c.SendExtensionResponse("mmu", "-1", "0", "1", strconv.Itoa(cash), strconv.Itoa(gold))
	case false:
		c.SendExtensionResponse("mmu", "-1", "0", "0", strconv.Itoa(-cash), strconv.Itoa(-gold))
	}

	return nil
}

func MuffinGameValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	cash, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	gold, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if c.Player.GetLevel() < 10 {
		return "Player not yet unlocked the feature", NOT_DECLARED
	}

	if event.GetEvent() == 3 {
		return "Cant play muffin minigame on winter!", NOT_DECLARED
	}

	if gold > 10000 || cash > 9999999 {
		return "Player cant bet more gold or cash", NOT_DECLARED
	}

	if gold < 0 || cash < 100 {
		return "Player bet lower gold or cash", NOT_DECLARED
	}

	// Just dont enable this
	if c.Player.GetCash() < cash || c.Player.GetGold() < gold {
		return "Player not have enough money", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
