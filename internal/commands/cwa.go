package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// cwa - C2S_CAFE_WALK
func CafeWalk(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	posX := req.Args[2]
	posY := req.Args[3]

	posXInt, err := strconv.Atoi(posX)
	if err != nil {
		return err
	}
	posYInt, err := strconv.Atoi(posY)
	if err != nil {
		return err
	}

	c.Player.Position = []int{posXInt, posYInt}

	c.Location.Broadcast("cwa", "-1", "0", strconv.Itoa(c.Player.ID), "0", posX, posY)

	return nil
}
