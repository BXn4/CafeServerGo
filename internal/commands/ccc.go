package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ccc - C2S_CAFE_COOK
func StartCooking(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CCC while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	var usingFancy int = 0

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return fmt.Errorf("Cant parse posX to int: %w", err)
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return fmt.Errorf("Cant parse posY to int: %w", err)
	}

	dishID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return fmt.Errorf("Cant parse dishID to int: %w", err)
	}

	isPrepared, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return fmt.Errorf("Cant parse isPrepared to int: %w", err)
	}

	usingFancy, err = strconv.Atoi(req.Args[6])
	if err != nil {

		return fmt.Errorf("Cant parse usingFancy to int: %w", err)
	}

	cookingTime := c.Player.GetDishMasteryDuration(dishID)

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if stove == nil {
		return fmt.Errorf("No stove found at: %v:%v", posX, posY)
	}

	if isPrepared == 0 {
		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return fmt.Errorf("Invalid ingredient ID: %w", err)
		}

		if gm.GetEvent() < dishInfo.Events {
			return fmt.Errorf("Invalid dish ID:, because theres no holiday %w", err)

		}

		ingredientsStr := dishInfo.Requirements
		ingredientsMap := make(map[int]int)
		ingredients := strings.Split(ingredientsStr, "#")
		for _, ingredient := range ingredients {
			parts := strings.Split(ingredient, "+")

			ingredientID, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("Error converting ingredient ID: %w", err)
			}

			ingredientAmount, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("Error converting amount: %w", err)
			}

			if ingredientID >= 1401 && usingFancy == 1 {
				ingredientsMap[ingredientID] = ingredientAmount
			} else if ingredientID < 1401 && usingFancy == 0 {
				ingredientsMap[ingredientID] = ingredientAmount
			}
		}

		for ingredientID, ingredientAmount := range ingredientsMap {
			fridgeInventoryAmount := c.Location.Cafe().GetFridgeInventory()[ingredientID]
			if ingredientAmount > fridgeInventoryAmount {
				return fmt.Errorf("Player not have enough ingredient in the fridge")
			}

			c.Location.Cafe().RemoveFromFridge(ingredientID, ingredientAmount)
		}

		stove.SetDishID(dishID)
		stove.SetFancyIng(usingFancy != 0)

		if usingFancy != 0 {
			c.Player.UpdateAchivementFancyCount()
		}

		// sweets = 0
		// meals = 1
		// soup = 2
		// salad = 3
		// vega = 4
		// snacks = 5

		dishCategory := dishInfo.DishCategory

		switch dishCategory {
		case 0:
			c.Player.UpdateAchivementServingCountSweets()
		case 1:
			c.Player.UpdateAchivementServingCountMeals()
		case 2:
			c.Player.UpdateAchivementServingCountSoups()
		case 3:
			c.Player.UpdateAchivementServingCountSalads()
		case 4:
			c.Player.UpdateAchivementServingCountVegans()
		case 5:
			c.Player.UpdateAchivementServingCountSnacks()
		}

	} else {
		stove.SetDishID(dishID)
		currentTime := time.Now().UTC()
		stove.SetStartedAt(&currentTime)
		finishesAt := stove.GetStartedAt().Add(time.Duration(cookingTime) * time.Second)
		stove.SetFinishesAt(&finishesAt)
	}

	c.Location.Broadcast(
		"ccc", "-1", "0",
		req.Args[2],
		req.Args[3],
		strconv.Itoa(dishID),
		strconv.Itoa(isPrepared),
		strconv.Itoa(usingFancy),
		strconv.Itoa(int(cookingTime)),
	)

	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())

	return nil
}
