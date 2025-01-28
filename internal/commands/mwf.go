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
)

func WheelOfFortune(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// If you have too much gifts return
	if len(c.Player.Gifts) >= 99 {
		c.SendExtensionResponse("mwf", "-1", "92")
		return nil
	}

	// If you already played the wheel
	if c.Player.PlayedWheel {
		c.SendExtensionResponse("mwf", "-1", "4")
	}

	reward := rand.Intn(16)
	rewardStr := strconv.Itoa(reward)

	switch reward {
	case 0:
		c.Player.Gold += 10
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1902+10")
	case 1:
		fallthrough
	case 9:
		amount := 300 + rand.Intn((401-300)/5)*5
		c.Player.EarnedChips(amount)
		amountStr := strconv.Itoa(amount)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1901+"+amountStr)
	case 2:
		fallthrough
	case 13:
		c.Player.AddGift(1450, 1, -1)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1450+1")
	case 3:
		decorations, err := utils.GetItems("deco")
		if err != nil {
			return err
		}
		choice := rand.Intn(len(decorations))
		dec := decorations[choice]
		idStr := strconv.Itoa(dec.ID)

		// Add won decoration to inventory
		mycafe, err := c.DB.GetCafeByPlayerID(c.Player.ID)
		if err != nil {
			return nil
		}
		mycafe.AddFurnitures(dec.ID, 1)

		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, idStr+"+1")
	case 4:
		fallthrough
	case 11:
		val := rand.Intn(int(99 - len(c.Player.Gifts)))
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
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, strings.Join(fanciesStr, "#"))
	case 5:
		fallthrough
	case 7:
		fallthrough
	case 15:
		val := rand.Intn(int(99 - len(c.Player.Gifts)))
		dishCount := utils.If(val > 3, 3, val)
		dishes, err := utils.GetItems("dish")
		if err != nil {
			return err
		}

		dishesStr := []string{}
		for i := 0; i < dishCount; i++ {
			// Get fancy and amount
			choice := rand.Intn(len(dishes))
			dish := dishes[choice]
			amount := rand.Intn(10) + 1

			// Add won dish to gifts
			c.Player.AddGift(dish.ID, amount, -1)

			// Add to won fancy list
			dishStr := fmt.Sprintf("%v+%v", dish.ID, amount)
			dishesStr = append(dishesStr, dishStr)
		}
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, strings.Join(dishesStr, "#"))
	case 6:
		fallthrough
	case 14:
		amount := (rand.Intn(150-60+1) + 60) * c.Player.GetLevel()
		c.Player.XP += amount
		amountStr := strconv.Itoa(amount)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1905+"+amountStr)
	case 8:
		amount := rand.Intn(5) + 1
		amountStr := strconv.Itoa(amount)
		c.Player.Gold += amount
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, "1902+"+amountStr)
	case 10:
		val := rand.Intn(int(99 - len(c.Player.Gifts)))
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
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, strings.Join(ingredientsStr, "#"))

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
		mycafe, err := c.DB.GetCafeByPlayerID(c.Player.ID)
		if err != nil {
			return nil
		}
		mycafe.AddFurnitures(decor.ID, 1)

		// Send msg
		idStr := strconv.Itoa(decor.ID)
		c.SendExtensionResponse("mwf", "-1", "0", rewardStr, idStr+"+1")
	}

	return nil
}
