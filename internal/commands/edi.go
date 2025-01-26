package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"

	"github.com/charmbracelet/log"
)

// edit - C2S_EDITOR_MODE
func EditorMode(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	status, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	// Do not allow to enter editor mode multiple times
	// Not casuing any issue yet
	if c.Location.Cafe().InEditorMode() && status != 0 {
		return nil
	}
	switch status {
	case 0:
		// If the player places an object where he stayed
		if c.Location.Cafe().GetObjectByPos(c.Player.Position[0], c.Player.Position[1]) != nil {
			for _, object := range c.Location.Cafe().GetObjects() {
				if object.IsDoor() {
					newStartPos := [2]int{
						utils.If(object.GetPos()[0] == 0, 1, object.GetPos()[0]),
						utils.If(object.GetPos()[1] == 0, 1, object.GetPos()[1]),
					}
					c.Location.Cafe().SetPlayerStart(newStartPos)
					c.Player.Position = newStartPos
					break
				}
			}
		}
		c.Location.Cafe().SetInEditorMode(false)
		c.Location.SetRunning(true)
		log.Printf("SET RUNING TRUE")
		c.SendExtensionResponse("edi", "-1", "0",
			strconv.Itoa(status), strconv.Itoa(c.Player.Position[0]), strconv.Itoa(c.Player.Position[1]), "") // <- Objects
	case 1:
		c.SendExtensionResponse("edi", "-1", "0",
			strconv.Itoa(status), strconv.Itoa(c.Player.Position[0]), strconv.Itoa(c.Player.Position[1]), "") // <- Objects
		c.Location.Cafe().SetInEditorMode(true)
		c.Location.Cafe().SetInEditorMode(true)
		c.Location.SetRunning(false)
		log.Printf("SET RUNING FALSE")
	default:
		return nil
	}
	return nil
}
