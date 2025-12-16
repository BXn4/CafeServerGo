package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_EDITOR_STORE_OBJECT,
		CommandConfig{
			Name:       "StoreObject",
			Identifier: responses.S2C_EDITOR_STORE_OBJECT,
			MinArgs:    4,
			MaxArgs:    4,
		},
		StoreObjectValidator,
		StoreObject,
	)
}

// est - C2S_EDITOR_STORE_OBJECT
func StoreObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, _ := strconv.Atoi(req.Args[2])
	objY, _ := strconv.Atoi(req.Args[3])

	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	c.Location.Cafe().AddFurnitures(int(obj.GetKind()), 1)

	objectInfo, _ := utils.GetItem(int(obj.GetKind()))
	c.Location.Cafe().RemoveObject(obj.GetPos())
	c.Location.Cafe().AddLuxury(-(objectInfo.Cash / 4000) + (objectInfo.Gold * 2))

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateLuxury(c.Location.Cafe().ID, c.Location.Cafe().GetLuxury())
	c.DB.UpdateFurnitureInventory(c.Location.Cafe().ID, c.Location.Cafe().FurnitureInventory.String())

	c.SendExtensionResponse("est", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(int(obj.GetKind())))
	return nil
}

func StoreObjectValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
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
	objX, objY := items[0], items[1]

	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	if obj == nil || obj.IsDoor() {
		return "No object found at the pos!", ERROR_EDITOR_STORE_FIELD_EMPTY
	}

	_, err = utils.GetItem(int(obj.GetKind()))
	if err != nil {
		return "Object info not found!", ERROR_EDITOR_WATCH_OUT
	}

	if obj.IsStove() || obj.IsCounter() {
		if obj.GetDishID() != -1 {
			return "Cant store object, because theres a dish on top!", ERROR_EDITOR_STORE_OBJ_DISH_ON_TOP_THIS
		}

		var stovesCount int = 0
		var countersCount int = 0

		for _, object := range c.Location.Cafe().Objects {
			if object.IsStove() {
				stovesCount++
			}
			if object.IsCounter() {
				countersCount++
			}
		}

		if stovesCount == 1 || countersCount == 1 {
			return "Cant store object, because the cafe needs 1 stove and 1 counter!", ERROR_EDITOR_WATCH_OUT
		}
	}

	if obj.IsFridge() || obj.IsDoor() {
		var fridgesCount int = 0
		var doorCount int = 0

		for _, object := range c.Location.Cafe().Objects {
			if object.IsDoor() {
				doorCount++
			}

			if object.IsFridge() {
				fridgesCount++
			}
		}

		if obj.IsFridge() {
			// (2-1)*50 <= 70
			if 50*(fridgesCount-1) <= c.Location.Cafe().GetFridgeCapacity() {
				return "Cant store the fridge, because the player got too many ingredietns!", ERROR_EDITOR_WATCH_OUT
			}
		}

		if c.Location.Cafe().GetFridgeFreeSpace() < 50*fridgesCount || doorCount == 1 {
			return "Cant store object, because the cafe needs 1 fridge and 1 door!", ERROR_EDITOR_WATCH_OUT
		}
	}

	return "Command ran without any errors.", NO_ERROR
}
