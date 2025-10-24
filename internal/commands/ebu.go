package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/simple"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// ebu - C2S_EDITOR_BUY_OBJECT
// TODO: Need to check level.
func BuyObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return err
	}
	objX, objY, objID, objRotation := items[0], items[1], items[2], items[3]

	// Dont allow players to modify the packet and sending us EBU while not in editor.
	if c.Location.IsRunning() {
		// Need to send the ID, because the client parse it / these.
		c.SendExtensionResponse("ebu", "-1", "38", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
		return nil
	}
	// Theres a object in that space. The game removes the door render on door drag. We need to enable that to place back to the og. pos,
	if obj := c.Location.Cafe().GetObjectByPosXY(objX, objY); obj != nil {
		// If not a door, give error.
		// I'm not sure if its really works, but if its works, we need to handle it.
		if !obj.IsDoor() {
			c.SendExtensionResponse("ebu", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil
		}
	}

	objectInfo, err := utils.GetItem(objID)
	if err != nil {
		return fmt.Errorf("Invalid object ID: %w", err)
	}

	// If the player does not have the object in their inventory, dont remove cash, gold
	if c.Location.Cafe().GetFurnitureInventory()[objID] == 0 {
		if objectInfo.Cash > c.Player.GetCash() {
			c.SendExtensionResponse("ebu", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil

		} else if objectInfo.Cash > 0 {
			c.Player.AddCash(-objectInfo.Cash)
		}

		if objectInfo.Gold > c.Player.GetGold() {
			c.SendExtensionResponse("ebu", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil
		} else if objectInfo.Gold > 0 {
			c.Player.AddGold(-objectInfo.Gold)
		}
	}

	// Need to add back the old wall in the inventory
	switch objectInfo.Group {
	case "Wall":
		oldWallID := c.Location.Cafe().GetTiles()[objX][objY]
		// If the old wall have luxury value, remove it from the Cafe
		c.Location.Cafe().AddLuxury(-(objectInfo.Cash / 4000) + (objectInfo.Gold * 2))
		c.Location.Cafe().AddFurnitures(oldWallID, 1)

		// Add the new wall
		c.Location.Cafe().SetTile(objX, objY, objID)

	case "Door":
		oldDoorPos := simple.NewPosition(
			utils.If(c.Location.Cafe().GetPlayerStart().X == 1, 0, c.Location.Cafe().GetPlayerStart().X),
			utils.If(c.Location.Cafe().GetPlayerStart().Y == 1, 0, c.Location.Cafe().GetPlayerStart().Y),
		)
		oldDoor := c.Location.Cafe().GetObjectByPos(oldDoorPos)
		// If the old door have luxury value, remove it from the Cafe
		obj, err := utils.GetDoor(int(oldDoor.GetKind()))
		if err != nil {
			return nil
		}
		c.Location.Cafe().AddLuxury(-(obj.Cash / 4000) + (obj.Gold * 2))
		// KIND = ID!!!
		c.Location.Cafe().GetFurnitureInventory()[int(oldDoor.GetKind())] = 1
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		c.Location.Cafe().RemoveObject(oldDoorPos)
	case "Deco":
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)

		c.Player.UpdateAchivementBoughtDecoration()

		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
	case "Fridge":
		numberOfFridges := c.Location.Cafe().GetFridgeMaxCapacity() / 50

		if numberOfFridges < utils.GetLevelFridgesLimit(c.Player.GetLevel()) {
			c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		} else {
			c.SendExtensionResponse("ebu", "-1", "3", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil // dont allow players to purchase more
		}
	case "Stove":
		numberOfStoves := 0

		for _, obj := range c.Location.Cafe().Objects {
			if obj.IsStove() {
				numberOfStoves++
			}
		}

		if numberOfStoves < utils.GetLevelStovesLimit(c.Player.GetLevel()) {
			c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		} else {
			c.SendExtensionResponse("ebu", "-1", "3", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil // dont allow players to purchase more
		}
	case "Counter":
		numberOfCounters := 0

		for _, obj := range c.Location.Cafe().Objects {
			if obj.IsCounter() {
				numberOfCounters++
			}
		}

		if numberOfCounters < utils.GetLevelStovesLimit(c.Player.GetLevel()) {
			c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		} else {
			c.SendExtensionResponse("ebu", "-1", "3", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil // dont allow players to purchase more
		}

	default:
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
	}
	// Works?
	// 0 / 4000 + 0 * 2 = 0 (if not cost cash) (if not cost gold)
	c.Location.Cafe().AddLuxury((objectInfo.Cash / 4000) + (objectInfo.Gold * 2))

	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateLuxury(c.Location.Cafe().ID, c.Location.Cafe().GetLuxury())

	c.SendExtensionResponse("ebu", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
	return nil
}
