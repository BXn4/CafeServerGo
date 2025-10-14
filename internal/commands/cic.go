package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
	"time"
)

// cic - S2C_CAFE_INSTANTCOOK
// TODO: level ratio
func InstantCook(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	// Dont allow players to modify the packet and sending us CIC while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	if c.Player.GetGold() < 1 /* && c.Player.GetInstantCookings < 1 */ {
		c.SendExtensionResponse("crc", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}

	stove := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	if stove == nil {
		return nil
	}

	c.Player.AddGold(-1)

	currentTime := time.Now().UTC()
	stove.SetStartedAt(&currentTime)
	stove.SetFinishesAt(&currentTime)

	c.Player.UpdateAchivementInstantCount()
	c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())

	c.Location.Broadcast("cic", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))

	return nil
}
