/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package job

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

/*func init() {
	RegisterCommand(requests.C2S_JOB_USER_ACTION,
		CommandConfig{
			Name:       "JobAction",
			Identifier: responses.S2C_JOB_USER_ACTION,
			MinArgs:    0,
			MaxArgs:    0,
		},
		JobActionValidator,
		JobAction,
	)
} */

var (
	JobCleanStart   = 0
	JobCleaned      = 1
	JobDeliverStart = 2
	JobDelivered    = 3
	JobPickupDish   = 4
	JobPickDownDish = 5
	JobMove         = 6

	EMPTY_HAND     = 62
	NOT_EMPTY_HAND = 63
)

func JobAction(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	/* if !c.Player.GetJobStarted() {
		// not started the job
		return nil
	}

	action, err := strconv.Atoi(req.Args[2])
	posX, err := strconv.Atoi(req.Args[3])
	posY, err := strconv.Atoi(req.Args[4])

	if err != nil {
		return fmt.Errorf("Failed to convert string to int: %s", err)
	}

	switch action {
	case JobCleanStart:
		if c.Player.Job.DishID != 0 {
			c.SendExtensionResponse("wua", strconv.Itoa(NOT_EMPTY_HAND))
			return nil
		}

		chair := c.Player.Job.Location.GetObjectByPosXY(posX, posY)

		if chair == nil {
			return fmt.Errorf("Invalid pos")
		}

		// clear it earlier, because if a new player joins the cafe, we need to send cleaned chairs, because only the player (waiter) can remove the dish visually
		// so if we want to send cleaned when a player joins the cafe, we need to send that player to clean the table every time (out of sync)

		tempDishID := chair.GetDishID()

		chair.SetDishID(-1)
		chair.SetDishStatus(-1)

		if c.Location.ReserveObject(chair) {
			start := agents.NewCafePoint(c.Player.Position, c.Player.Job.Location)
			end := agents.NewCafePoint(chair.GetPos(), c.Player.Job.Location)
			_, distance, found := agents.Path(start, end)

			if !found {
				return fmt.Errorf("Path not found!")
			}

			c.Location.Broadcast("wua", strconv.Itoa(JobCleanStart), strconv.Itoa(c.Player.ID), strconv.Itoa(chair.GetPos().X), strconv.Itoa(chair.GetPos().Y))

			if !c.Location.TryStepSleep((time.Duration(distance-1) * 550) * time.Millisecond) {
				// cancelled or cafe was set to stop

				// to allow other waiters to clean it
				chair.SetDishID(tempDishID)
				chair.SetDishStatus(3)

				c.Location.UnreserveObject(chair)
				return nil
			}
		}

	}
	*/
	return nil
}
