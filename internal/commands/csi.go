package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/object"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

func StoveDeliverInfo(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CSD while in editor.
	if !c.Location.IsRunning() {
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

	stove := c.Location.Cafe().GetObjectByPosXY(stoveX, stoveY)

	if !stove.GetIsRotten() {

		// Choose counter that is empty or has the same food type
		var counter *object.Object
		for _, object := range c.Location.Cafe().GetObjects() {
			if !object.IsCounter() {
				continue
			}

			if stove.GetDishID() == object.GetDishID() {
				counter = object
				break
			} else if object.GetDishID() == -1 {
				counter = object
			}
		}

		// Set args
		counterX := utils.If(counter != nil, counter.GetPos().X, -1)
		counterY := utils.If(counter != nil, counter.GetPos().Y, -1)
		status := utils.If(counter != nil, "0", "37")

		/*
			py	%xt%csi%-1%0%1%5%3%6%
			go	%xt%csi%-1%0%0%1%5%3%6%
		*/
		c.Location.Broadcast(
			"csi", "-1",
			status,
			req.Args[2],
			req.Args[3],
			strconv.Itoa(counterX),
			strconv.Itoa(counterY),
		)
	}

	return nil
}
