package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
)

// ese - C2S_EDITOR_SELL_OBJECT
func SellObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
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

	sellAmount, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return err
	}

	// Dont allow players to modify the packet and sending us ESE while not in editor.
	if c.Location.IsRunning() {
		c.SendExtensionResponse("ese", "-1", "38", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}

	// If the object id is not valid
	objectInfo, err := utils.GetItem(objID)
	if err != nil {
		return fmt.Errorf("Invalid object ID: %w", err)
	}

	if objX != -1 && objY != -1 {
		obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)
		if obj == nil {
			c.SendExtensionResponse("ese", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY))
			return nil
		}

		// If the player is not sending from inventory, remove luxury
		c.Location.Cafe().AddLuxury(-(objectInfo.Cash / 4000) + (objectInfo.Gold * 2))

		c.Location.Cafe().RemoveObject(obj.GetPos())
	} else {
		// If the player wants to send more than they have
		if c.Location.Cafe().GetFurnitureInventory()[objID] < sellAmount {
			c.SendExtensionResponse("ese", "-1", "32", strconv.Itoa(objX), strconv.Itoa(objY))
			return nil
		}

		c.Location.Cafe().RemoveFurnitures(objID, sellAmount)

	}

	c.Player.AddCash(sellAmount * int(math.Round(float64(objectInfo.Cash)*0.2+float64(objectInfo.Gold)*0.2)))

	c.Player.UpdateAchivementSoldItems()

	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateLuxury(c.Location.Cafe().ID, c.Location.Cafe().GetLuxury())

	c.SendExtensionResponse("ese", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objID), strconv.Itoa(sellAmount))
	return nil
}
