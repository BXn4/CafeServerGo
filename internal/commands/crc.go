package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
	"time"
)

// crc - Recook
func Recook(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// Dont allow players to modify the packet and sending us CRC while in editor.
	if !c.Location.IsRunning() {
		return nil
	}

	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	if c.Player.GetGold() < 1 {
		c.SendExtensionResponse("crc", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}

	stove := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	if stove == nil {
		return nil
	}

	if stove.GetIsRotten() { // if its really rotten
		c.Player.AddGold(-1)

		// The recook works different!!
		// this.rottenCookTime = getTimer() + currentDish.baseDuration * 60000 / CafeConstants.timeFactor;
		// https://www.youtube.com/watch?v=GH_Fw6yAjJo
		cookingTime := c.Player.GetDishMasteryDuration(stove.GetDishID()) // returns in seconds
		currentTime := time.Now().UTC()

		if cookingTime < 60*60 {
			startedAt := currentTime.Add(-time.Duration(cookingTime)*time.Second - time.Hour) // from the current hour remove X seconds, then 1 hour
			finishedAt := currentTime.Add(time.Duration(cookingTime)*time.Second - time.Hour)

			stove.SetStartedAt(&startedAt)
			stove.SetFinishesAt(&finishedAt)
		} else {
			startedAt := currentTime
			finishedAt := currentTime

			stove.SetStartedAt(&startedAt)
			stove.SetFinishesAt(&finishedAt)
		}

		c.Player.UpdateAchivementOvercookedFoods()

		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

		c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())
		c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())

		c.Location.Broadcast("crc", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY))
	} else {
		c.SendExtensionResponse("crc", "-1", "4", strconv.Itoa(objX), strconv.Itoa(objY))
	}

	return nil
}
