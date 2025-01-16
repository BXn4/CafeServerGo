package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// emo - C2S_EDITOR_MOVE_OBJECT
func MoveObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	oldObjX, err := strconv.Atoi(req.Args[2])
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
	}
	// Dont allow players to modify the packet and sending us ERO while not in editor.
	if !c.Location.Cafe().InEditorMode {
		c.SendExtensionResponse("emo", "-1", "38", strconv.Itoa(oldObjX), strconv.Itoa(oldObjY))
		return nil
	}
	obj := c.Location.Cafe().GetObjectByPos(oldObjX, oldObjY)
	if obj == nil {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjX), strconv.Itoa(oldObjY))
		return nil
	}
	// Dont allow players to place in the player start
	if newObjX == c.Location.Cafe().GetPlayerStart()[0] && newObjY == c.Location.Cafe().GetPlayerStart()[1] {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjX), strconv.Itoa(oldObjY))
		return nil
	}
	if c.Location.Cafe().GetObjectByPos(newObjX, newObjY) != nil {
		c.SendExtensionResponse("emo", "-1", "39", strconv.Itoa(oldObjX), strconv.Itoa(oldObjY))
		return nil
	}
	c.Location.Cafe().RemoveObject(oldObjX, oldObjY)
	c.Location.Cafe().AddNewObject(newObjX, newObjY, int(obj.Kind), int(obj.Rotation))
	c.SendExtensionResponse("emo", "-1", "0", strconv.Itoa(oldObjX), strconv.Itoa(oldObjY), strconv.Itoa(newObjX), strconv.Itoa(newObjX))
	return nil
}
