package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/simple"
	"cafego/internal/types/requests"
	"strconv"
)

// cwa - C2S_CAFE_WALK
func CafeWalk(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CWA while in editor.
	if !c.Location.IsRunning() {
		return nil
	}
	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	if c.Location.Cafe().GetObjectByPosXY(posX, posY) != nil {
		c.SendExtensionResponse("cwa", "-1", "23")
		return nil
	}

	c.Player.Position = simple.NewPosition(posX, posY)

	c.Location.Broadcast("cwa", "-1", "0", strconv.Itoa(c.Player.ID), "0", req.Args[2], req.Args[3])

	return nil
}
