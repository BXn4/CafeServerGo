package cafe

import (
	"cafego/internal/client"
	"cafego/internal/commands"
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
	commands.RegisterCommand(requests.C2S_CAFE_COOK,
		commands.CommandConfig{
			Name:        "Cook",
			Identifier:  responses.S2C_CAFE_COOK,
			Description: "Starts cooking the dishes / prepare done handling",
			Args:        "{posX} {posY} {dishID} {isPrepared} {isUsingFancy} {cookingTimeInSeconds}",
			MinArgs:     7,
			MaxArgs:     7,
		},
		StartCookingValidator,
		StartCooking,
		StartCookingDBSaver,
	)
}

// ccc - C2S_CAFE_COOK
func StartCooking(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	var usingFancy int = 0

	posX, _ := strconv.Atoi(req.Args[2])
	posY, _ := strconv.Atoi(req.Args[3])
	dishID, _ := strconv.Atoi(req.Args[4])
	isPrepared, _ := strconv.Atoi(req.Args[5])
	usingFancy, _ = strconv.Atoi(req.Args[6])

	println("1")

	cookingTime := c.Player.GetDishMasteryDuration(dishID)
	println("2")
	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)
	println("3")

	if isPrepared == 0 {
		dishInfo, _ := utils.GetDish(dishID)
		println("4")

		ingredientsStr := dishInfo.Requirements
		ingredientsMap := make(map[int]int)
		ingredients := strings.Split(ingredientsStr, "#")
		for _, ingredient := range ingredients {
			parts := strings.Split(ingredient, "+")
			ingredientID, _ := strconv.Atoi(parts[0])
			ingredientAmount, _ := strconv.Atoi(parts[1])

			if usingFancy == 1 {
				if ingredientID >= 1401 {
					ingredientsMap[ingredientID] = ingredientAmount
				}
			} else {
				if ingredientID < 1401 {
					ingredientsMap[ingredientID] = ingredientAmount
				}
			}
		}

		println("5")

		for ingredientID, ingredientAmount := range ingredientsMap {
			c.Location.Cafe().RemoveFromFridge(ingredientID, ingredientAmount)
		}

		stove.SetDishID(dishID)
		stove.SetFancyIng(usingFancy != 0)

		if usingFancy == 1 {
			c.Player.UpdateAchivementFancyCount()
		}

		println("6")

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
		cm.Identifier, "-1", "0",
		strconv.Itoa(posX),
		strconv.Itoa(posY),
		strconv.Itoa(dishID),
		strconv.Itoa(isPrepared),
		strconv.Itoa(usingFancy),
		strconv.Itoa(int(cookingTime)),
	)

	return nil
}

func StartCookingValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us CCC while in editor.
	if !c.Location.IsRunning() {
		return "The location is not running", commands.LOCATION_NOT_RUNNING
	}

	if c.Location.Cafe().GetOwnerID() != c.Player.GetID() {
		return "Not the owner!", commands.NOT_DECLARED
	}

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	dishID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if !c.Player.GetIsTutorialCompleted() {
		if posX != 1 && posY != 7 && dishID != 1201 {
			return "Invalid tutorial values, needed: 1:7 for pos, and 1201 for dish!", commands.INVALID_VALUE
		}
	}

	isPrepared, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if req.Args[5] != "0" && req.Args[5] != "1" {
		return fmt.Sprintf("Invalid args for boolean: %v", req.Args[5]), commands.INVALID_ARGS
	}

	usingFancy, err := strconv.Atoi(req.Args[6])
	if err != nil {
		return "Cant convert string to int!", commands.CONVERT_ERROR
	}

	if req.Args[6] != "0" && req.Args[6] != "1" {
		return fmt.Sprintf("Invalid args for boolean: %v", req.Args[6]), commands.INVALID_ARGS
	}

	stove := c.Location.Cafe().GetObjectByPosXY(posX, posY)
	if stove == nil {
		return fmt.Sprintf("No stove found at: %d:%d", posX, posY), commands.NOT_DECLARED
	}

	if usingFancy == 1 {
		if c.Player.GetLevel() < 9 {
			return "Cant use fancy, because the player level is not enough!", commands.NOT_DECLARED
		}
	}

	if isPrepared == 0 && stove.GetDishID() != -1 {
		return "Cant start cook on that stove, because the stove already have a dish!", commands.NOT_DECLARED
	}

	if isPrepared == 0 {
		dishInfo, err := utils.GetDish(dishID)
		if err != nil {
			return "Invalid dish ID", commands.INVALID_VALUE
		}

		if event.GetEvent() < dishInfo.Events {
			return "Invalid dish ID, because theres no holiday event.", commands.NOT_DECLARED
		}

		if c.Player.GetLevel() < dishInfo.Level {
			return "Player level is not enough to cook the dish", commands.NOT_DECLARED
		}

		ingredientsStr := dishInfo.Requirements
		ingredientsMap := make(map[int]int)
		ingredients := strings.Split(ingredientsStr, "#")
		for _, ingredient := range ingredients {
			parts := strings.Split(ingredient, "+")

			ingredientID, err := strconv.Atoi(parts[0])
			if err != nil {
				return "Cant convert string to int!", commands.CONVERT_ERROR
			}

			ingredientAmount, err := strconv.Atoi(parts[1])
			if err != nil {
				return "Cant convert string to int!", commands.CONVERT_ERROR
			}

			if usingFancy == 1 {
				if ingredientID >= 1401 {
					ingredientsMap[ingredientID] = ingredientAmount
				}
			} else {
				if ingredientID < 1401 {
					ingredientsMap[ingredientID] = ingredientAmount
				}
			}
		}

		for ingredientID, ingredientAmount := range ingredientsMap {
			fridgeInventoryAmount := c.Location.Cafe().GetFridgeInventory()[ingredientID]
			if ingredientAmount > fridgeInventoryAmount {
				return "Player not have enough ingredient in the fridge", commands.NOT_DECLARED
			}
		}
	}

	return "Command ran without any errors.", commands.NO_ERROR
}

func StartCookingDBSaver(c *client.Client) error {
	if !c.Player.GetIsTutorialCompleted() {
		return nil
	}

	c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())
	c.DB.UpdateFridgeInventory(c.Location.Cafe().GetID(), c.Location.Cafe().GetFridgeInventory().String())

	return nil
}
