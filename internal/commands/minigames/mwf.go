package minigames

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_WHEELOFFORTUNE,
		commands.CommandConfig{
			Name:         "WheelOfFortuneMinigame",
			Description:  "Fortune wheel minigame",
			Identifier:   responses.S2C_WHEELOFFORTUNE,
			Args:         "{winStr}",
			MinArgs:      0,
			MaxArgs:      0,
			FeatureLevel: 9,
		},
		WheelOfFortuneValidator,
		WheelOfFortune,
		WheelOfFortuneDBSaver,
	)
}

// min level 9

func WheelOfFortune(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	if c.Player.GetPlayedWheel() {
		c.Player.AddGold(-1)
	} else {
		c.Player.SetPlayedWheel(true)
	}

	c.Player.UpdateAchivementWheelOfFortune()

	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
	c.DB.UpdatePlayedWheel(c.Player.GetID(), c.Player.GetPlayedWheel())

	/* random := rand.New(rand.NewSource(time.Now().UnixNano()))

	rewards := []int{
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		9, 13, 3, 11, 15, 14, 10, 12,
		0, 8, // gold rewards
	}

	reward := rewards[random.Intn(len(rewards))] */
	reward := 10
	rewardStr := strconv.Itoa(reward)

	switch reward {
	case 0:
		c.Player.AddGold(10)
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, "1902+10")
	case 9:
		cashAmount := rand.Intn((401 - 300) / 5)
		factor := float64(c.Player.GetLevel()) * 0.5
		amount := 300 + int(float64(cashAmount)*factor)
		c.Player.AddCash(amount)
		amountStr := strconv.Itoa(amount)
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, "1901+"+amountStr)
	case 13:
		c.Player.AddGift(1450, 1, -1)
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, "1450+1")
	case 3:
		basicDecorIDs := make([]int, 9)
		for i := 0; i < 9; i++ {
			basicDecorIDs[i] = 501 + i
		}

		wonDecorID := basicDecorIDs[rand.Intn(len(basicDecorIDs))]
		idStr := strconv.Itoa(wonDecorID)

		// Add won decoration to inventory
		mycafe, err := c.DB.GetCafeByPlayerID(c.Player.GetID())
		if err != nil {
			return nil
		}
		mycafe.AddFurnitures(wonDecorID, 1)

		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, idStr+"+1")
	case 11:
		val := rand.Intn(int(99 - len(c.Player.GetGifts())))
		fancyCount := utils.If(val > 3, 3, val)
		fancies, err := utils.GetItems("fancy")
		if err != nil {
			return err
		}

		fanciesStr := []string{}
		for i := 0; i < fancyCount; i++ {
			// Get fancy and amount
			choice := rand.Intn(len(fancies))
			fancy := fancies[choice]
			amount := rand.Intn(10) + 1

			// Add won fancy to gifts
			c.Player.AddGift(fancy.ID, amount, -1)

			// Add to won fancy list
			fancyStr := fmt.Sprintf("%v+%v", fancy.ID, amount)
			fanciesStr = append(fanciesStr, fancyStr)
		}
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, strings.Join(fanciesStr, "#"))
	case 15:
		val := rand.Intn(int(99 - len(c.Player.GetGifts())))
		dishCount := utils.If(val > 3, 3, val)
		dishes, err := utils.GetItems("dish")
		if err != nil {
			return err
		}

		var validDishes []utils.Wod

		for _, dish := range dishes {
			if event.GetEvent() <= dish.Events {
				validDishes = append(validDishes, dish)
			}
		}

		dishesStr := []string{}
		for i := 0; i < dishCount; i++ {
			// Get fancy and amount
			choice := rand.Intn(len(validDishes))
			dish := dishes[choice]
			amount := rand.Intn(10) + 1

			// Add won dish to gifts
			c.Player.AddGift(dish.ID, amount, -1)

			// Add to won fancy list
			dishStr := fmt.Sprintf("%v+%v", dish.ID, amount)
			dishesStr = append(dishesStr, dishStr)
		}
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, strings.Join(dishesStr, "#"))
	case 14:
		amount := (rand.Intn(150-60+1) + 60) * c.Player.GetLevel()
		c.Player.AddXP(amount)
		amountStr := strconv.Itoa(amount)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1905+"+amountStr)
	case 8:
		amount := rand.Intn(5) + 1
		amountStr := strconv.Itoa(amount)
		c.Player.AddGold(amount)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1902+"+amountStr)
	case 10:
		val := rand.Intn(int(99 - len(c.Player.GetGifts())))
		ingredientCount := utils.If(val > 3, 3, val)
		ingredients, err := utils.GetItems("ingredient")
		if err != nil {
			return err
		}
		ingredientsStr := []string{}
		for i := 0; i < ingredientCount; i++ {
			// Get fancy and amount
			choice := rand.Intn(len(ingredients))
			ingredient := ingredients[choice]
			amount := rand.Intn(10) + 1

			// Add won ingredient to gifts
			c.Player.AddGift(ingredient.ID, amount, -1)

			// Add to won ingredient list
			dishStr := fmt.Sprintf("%v+%v", ingredient.ID, amount)
			ingredientsStr = append(ingredientsStr, dishStr)
		}
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, strings.Join(ingredientsStr, "#"))

	case 12:
		// Get all decoration
		decors, err := utils.GetItems("deco")
		if err != nil {
			return err
		}

		// Select won decor
		choice := rand.Intn(len(decors))
		decor := decors[choice]

		// Add won decoration to inventory
		mycafe, err := c.DB.GetCafeByPlayerID(c.Player.GetID())
		if err != nil {
			return nil
		}
		mycafe.AddFurnitures(decor.ID, 1)

		// Send msg
		idStr := strconv.Itoa(decor.ID)
		c.SendExtensionResponse(cm.Identifier, "-1", "0", rewardStr, idStr+"+1")
	}

	return nil
}

func WheelOfFortuneValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Player not yet unlocked the feature", commands.NOT_DECLARED
	}

	if c.Player.GetPlayedWheel() && c.Player.GetGold() <= 0 {
		return "Player not have enough money", commands.NOT_ENOUGHT_MONEY
	}

	// If you have too much gifts return
	if len(c.Player.GetGifts()) >= 99 {
		return "Player not have free gift space", commands.ERROR_WHEEL_MINIGAME_NO_GIFT_SPACE_LEFT
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func WheelOfFortuneDBSaver(c *client.Client) error {
	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateGold(c.Player.GetID(), c.Player.GetGold())
	c.DB.UpdateXP(c.Player.GetID(), c.Player.GetXP())
	c.DB.UpdateGifts(c.Player.GetID(), c.Player.GetGifts().String())
	c.DB.UpdateFurnitureInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFurnitureInventory().String())

	return nil
}
