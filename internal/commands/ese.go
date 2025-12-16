package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/balancing"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_EDITOR_SELL_OBJECT,
		CommandConfig{
			Name:       "SellObject",
			Identifier: responses.S2C_EDITOR_SELL_OBJECT,
			MinArgs:    6,
			MaxArgs:    6,
		},
		SellObjectValidator,
		SellObject,
	)
}

// ese - C2S_EDITOR_SELL_OBJECT
func SellObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])
	objID, _ := strconv.Atoi(req.Args[4])
	sellAmount, _ := strconv.Atoi(req.Args[5])

	objectInfo, _ := utils.GetItem(objID)
	// from inventory
	if objX == -1 && objY == -1 {
		c.Location.Cafe().RemoveFurnitures(objID, sellAmount)
	} else {
		obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)

		// If the player is not sending from inventory, remove luxury
		c.Location.Cafe().AddLuxury(-(objectInfo.Cash / 4000) + (objectInfo.Gold * 2))

		c.Location.Cafe().RemoveObject(obj.GetPos())
	}

	c.Player.AddCash(sellAmount * int(math.Round(float64(objectInfo.Cash)*(float64(balancing.BalancingConstants.SellFactorCash)/100)+float64(objectInfo.Gold)*(float64(balancing.BalancingConstants.SellFactorGold)/100))))

	c.Player.UpdateAchivementSoldItems(sellAmount)

	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateLuxury(c.Location.Cafe().ID, c.Location.Cafe().GetLuxury())

	c.SendExtensionResponse("ese", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(sellAmount))
	return nil
}

func SellObjectValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	// Dont allow players to modify the packet and sending us ESE while not in editor.
	if c.Location.IsRunning() {
		return "The location is running", ERROR_EDITOR_ONLY_IN_EDITOR
	}

	if c.Location.Cafe().GetPlayerID() != c.Player.ID {
		return "Not the owner!", NOT_DECLARED
	}

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return "Cant convert string to int!", CONVERT_ERROR
	}
	objX, objY, objID, sellAmount := items[0], items[1], items[2], items[3]

	_, err = utils.GetItem(objID)
	if err != nil {
		return "Object info not found!", ERROR_EDITOR_WATCH_OUT
	}

	// from inventory
	if objX == -1 && objY == -1 {
		if c.Location.Cafe().GetFurnitureInventory()[objID] < sellAmount {
			return "Object or amount not found in the inventory!", ERROR_EDITOR_SELL_ITEM_NOT_IN_STOCK
		}
	} else {
		if c.Location.Cafe().GetObjectByPosXY(objX, objY) == nil {
			return "No object found at the pos!", ERROR_EDITOR_POSITION_NOT_VALID
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
