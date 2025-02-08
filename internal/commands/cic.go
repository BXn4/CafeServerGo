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

	if c.Player.GetGold() < 1 {
		c.Location.Broadcast("cic", "-1", "4")
	}

	stove := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	if stove == nil {
		return nil
	}

	currentTime := time.Now().UTC()
	stove.SetStartedAt(&currentTime)
	stove.SetFinishesAt(&currentTime)

	c.Location.Broadcast("cic", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))

	return nil
}
