package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// min level = 9 premium dish

func init() {
	RegisterCommand(requests.C2S_CAFE_COOK,
		CommandConfig{
			Name:       "Cook",
			Identifier: responses.S2C_CAFE_COOK,
			MinArgs:    7,
			MaxArgs:    7,
		},
		StartCookingValidator,
		StartCooking,
	)
}

// ccc - C2S_CAFE_COOK
func StartCooking(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var usingFancy int = 0

	posX, _ := strconv.Atoi(req.Args[2])
	posY, _ := strconv.Atoi(req.Args[3])
	dishID, _ := strconv.Atoi(req.Args[4])
	isPrepared, _ := strconv.Atoi(req.Args[5])
	usingFancy, _ = strconv.Atoi(req.Args[6])

	cookingTime := c.Player.GetDishMasteryDuration(dishID)
	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)

	if isPrepared == 0 {
		dishInfo, _ := utils.GetDish(dishID)

		ingredientsStr := dishInfo.Requirements
		ingredientsMap := make(map[int]int)
		ingredients := strings.Split(ingredientsStr, "#")
		for _, ingredient := range ingredients {
			parts := strings.Split(ingredient, "+")
			ingredientID, _ := strconv.Atoi(parts[0])
			ingredientAmount, _ := strconv.Atoi(parts[1])
			if ingredientID >= 1401 && usingFancy == 1 {
				ingredientsMap[ingredientID] = ingredientAmount
			} else if ingredientID < 1401 && usingFancy == 0 {
				ingredientsMap[ingredientID] = ingredientAmount
			}
		}

		for ingredientID, ingredientAmount := range ingredientsMap {
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

	if c.Player.IsTutorialCompleted {
		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().Objects.StringForDB())
		c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().GetFridgeInventory().String())
	}

	return nil
}

func StartCookingValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CCC while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", LOCATION_NOT_RUNNING
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	dishID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if !c.Player.IsTutorialCompleted {
		if posX != 1 && posY != 7 && dishID != 1201 {
			return "Invalid tutorial values, needed: 1:7 for pos, and 1201 for dish!", INVALID_VALUE
		}
	}

	isPrepared, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if req.Args[5] != "0" && req.Args[5] != "1" {
		return fmt.Sprintf("Invalid args for boolean: %v", req.Args[5]), INVALID_ARGS
	}

	usingFancy, err := strconv.Atoi(req.Args[6])
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}

	if req.Args[6] != "0" && req.Args[6] != "1" {
		return fmt.Sprintf("Invalid args for boolean: %v", req.Args[6]), INVALID_ARGS
	}

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)
	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), NOT_DECLARED
	}

	if usingFancy == 1 {
		if c.Player.GetLevel() < 9 {
			return "Cant use fancy, because the player level is not enough!", NOT_DECLARED
		}
	}

	if isPrepared == 0 && stove.GetDishID() != -1 {
		return "Cant start cook on that stove, because the stove already have a dish!", NOT_DECLARED
	}

	if isPrepared == 0 {
		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Invalid dish ID", INVALID_VALUE
		}

		if event.GetEvent() < dishInfo.Events {
			return "Invalid dish ID, because theres no holiday event.", NOT_DECLARED
		}

		if c.Player.GetLevel() < dishInfo.Level {
			return "Player level is not enough to cook the dish", NOT_DECLARED
		}

		ingredientsStr := dishInfo.Requirements
		ingredientsMap := make(map[int]int)
		ingredients := strings.Split(ingredientsStr, "#")
		for _, ingredient := range ingredients {
			parts := strings.Split(ingredient, "+")

			ingredientID, err := strconv.Atoi(parts[0])
			if err != nil {
				return "Cant convert string to int!", CONVERT_ERROR
			}

			ingredientAmount, err := strconv.Atoi(parts[1])
			if err != nil {
				return "Cant convert string to int!", CONVERT_ERROR
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
				return "Player not have enough ingredient in the fridge", NOT_DECLARED
			}
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
