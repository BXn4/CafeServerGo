package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/object"
	"cafego/internal/types/requests"
	"strconv"
)

// ero - C2S_EDITOR_ROTATE_OBJECT
func RotateObject(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}
	objRotation, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}
	// Dont allow players to modify the packet and sending us ERO while not in editor.
	if c.Location.IsRunning() {
		c.SendExtensionResponse("ero", "-1", "38", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)
	if obj == nil {
		c.SendExtensionResponse("ero", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	if objRotation < 0 || objRotation > 3 {
		// Different rotation breaks the Cafe!!
		c.SendExtensionResponse("ero", "-1", "39", strconv.Itoa(objX), strconv.Itoa(objY))
		return nil
	}
	obj.SetRotation(object.CafeObjectRotation(objRotation))
	c.SendExtensionResponse("ero", "-1", "0", strconv.Itoa(objX), strconv.Itoa(objY), strconv.Itoa(objRotation))
	return nil
}
