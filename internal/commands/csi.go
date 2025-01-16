package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

func StoveDeliverInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CSD while in editor.
	if c.Location.Cafe().InEditorMode {
		return nil
	}

	stoveX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	stoveY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	stove := c.Location.Cafe().GetObjectByPos(stoveX, stoveY)

	// Choose counter that is empty or has the same food type
	var counter *objects.CafeObject
	for _, object := range c.Location.Cafe().Objects {
		if !object.IsCounter() {
			continue
		}

		if stove.DishID == object.DishID {
			counter = object
			break
		} else if object.DishID == -1 {
			counter = object
		}
	}

	// Set args
	counterX := utils.If(counter != nil, counter.Pos[0], -1)
	counterY := utils.If(counter != nil, counter.Pos[1], -1)
	status := utils.If(counter != nil, "0", "37")

	/*
		py	%xt%csi%-1%0%1%5%3%6%
		go	%xt%csi%-1%0%0%1%5%3%6%
	*/
	c.SendExtensionResponse(
		"csi", "-1",
		status,
		req.Args[2],
		req.Args[3],
		strconv.Itoa(counterX),
		strconv.Itoa(counterY),
	)
	return nil
}
