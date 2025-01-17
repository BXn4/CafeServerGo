package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

// ebu - C2S_EDITOR_BUY_OBJECT
// TODO: Need to check level.
func BuyObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	objID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}

	objRotation, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return err
	}

	// Dont allow players to modify the packet and sending us EBU while not in editor.
	if !c.Location.Cafe().InEditorMode {
		// Need to send the ID, because the client parse it / these.
		c.SendExtensionResponse("ebu", "-1", "38", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
		return nil
	}
	// Theres a object in that space. The game removes the door render on door drag. We need to enable that to place back to the og. pos,
	if obj := c.Location.Cafe().GetObjectByPos(objX, objY); obj != nil {
		// If not a door, give error.
		// I'm not sure if its really works, but if its works, we need to handle it.
		if !obj.IsDoor() {
			c.SendExtensionResponse("ebu", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil
		}
	}

	objectInfo, err := utils.GetObject(objID)
	if err != nil {
		return fmt.Errorf("Invalid object ID: %w", err)
	}

	// If the player does not have the object in their inventory, dont remove cash, gold
	if c.Location.Cafe().FurnitureInventory[objID] != 0 {
		if objectInfo.Cash != 0 && objectInfo.Cash > c.Player.Cash {
			// Need to send the ID, because the client parse it / these.
			c.SendExtensionResponse("ebu", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil
		}

		if objectInfo.Gold != 0 && objectInfo.Gold > c.Player.Gold {
			// Need to send the ID, because the client parse it / these.
			c.SendExtensionResponse("ebu", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
			return nil
		}

		c.Player.Cash -= objectInfo.Cash
		c.Player.Gold -= objectInfo.Gold
	}

	// Need to add back the old wall in the inventory
	if objectInfo.Group == "Wall" {
		oldWallID := c.Location.Cafe().Tiles[objX][objY]
		// If the old wall have luxury value, remove it from the Cafe
		c.Location.Cafe().Luxury -= (objectInfo.Cash / 4000) + (objectInfo.Gold * 2)
		if c.Location.Cafe().FurnitureInventory[oldWallID] != 0 {
			c.Location.Cafe().FurnitureInventory[oldWallID] += 1
		} else {
			c.Location.Cafe().FurnitureInventory[oldWallID] = 1
		}
		// Add the new wall
		c.Location.Cafe().Tiles[objX][objY] = objID

	} else if objectInfo.Group == "Door" {
		oldDoorPos := []int{
			utils.If(c.Location.Cafe().PlayerStart[0] == 1, 0, c.Location.Cafe().PlayerStart[0]),
			utils.If(c.Location.Cafe().PlayerStart[1] == 1, 0, c.Location.Cafe().PlayerStart[1]),
		}
		oldDoor := c.Location.Cafe().GetObjectByPos(oldDoorPos[0], oldDoorPos[1])
		// If the old door have luxury value, remove it from the Cafe
		obj, err := utils.GetDoor(int(oldDoor.Kind))
		if err != nil {
			return nil
		}
		c.Location.Cafe().Luxury -= (obj.Cash / 4000) + (obj.Gold * 2)
		// KIND = ID!!!
		c.Location.Cafe().FurnitureInventory[int(oldDoor.Kind)] = 1
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
		c.Location.Cafe().RemoveObject(oldDoorPos[0], oldDoorPos[1])
	} else {
		c.Location.Cafe().AddNewObject(objX, objY, objID, objRotation)
	}
	// Works?
	// 0 / 4000 + 0 * 2 = 0 (if not cost cash) (if not cost gold)
	c.Location.Cafe().Luxury += (objectInfo.Cash / 4000) + (objectInfo.Gold * 2)
	c.SendExtensionResponse("ebu", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(objRotation))
	return nil
}
