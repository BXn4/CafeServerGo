package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/simple"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

// emo - C2S_EDITOR_MOVE_OBJECT
func MoveObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return err
	}

	oldObjPos := simple.NewPosition(items[0], items[1])
	newObjPos := simple.NewPosition(items[2], items[3])

	/*oldObjX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	oldObjY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}
	newObjX, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}
	newObjY, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return err
		}*/
	// Dont allow players to modify the packet and sending us ERO while not in editor.
	if c.Location.IsRunning() {
		c.SendExtensionResponse("emo", "-1", "38", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y))
		return nil
	}
	obj := c.Location.Cafe().GetObjectByPos(oldObjPos)
	if obj == nil {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y))
		return nil
	}
	// Dont allow players to place in the player start
	if newObjPos == c.Location.Cafe().GetPlayerStart() {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y))
		return nil
	}
	if c.Location.Cafe().GetObjectByPos(newObjPos) != nil {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y))
		return nil
	}

	obj.SetPosXY(newObjPos.X, newObjPos.Y)

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())

	c.SendExtensionResponse("emo", "-1", "0", strconv.Itoa(oldObjPos.X), strconv.Itoa(oldObjPos.Y), strconv.Itoa(newObjPos.X), strconv.Itoa(newObjPos.X))
	return nil
}
