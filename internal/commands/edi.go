package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"log"
	"strconv"
)

// edit - C2S_EDITOR_MODE
func EditorMode(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	status, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}
	// Do not allow to enter editor mode multiple times
	// Not casuing any issue yet
	if c.Location.Cafe().InEditorMode && status != 0 {
		return nil
	}
	switch status {
	case 0:
		// If the player places an object where he stayed
		if c.Location.Cafe().GetObjectByPos(c.Player.Position[0], c.Player.Position[1]) != nil {
			for _, object := range c.Location.Cafe().Objects {
				if object.IsDoor() {
					newStartPos := []int{
						utils.If(object.Pos[0] == 0, 1, object.Pos[0]),
						utils.If(object.Pos[1] == 0, 1, object.Pos[1]),
					}
					c.Location.Cafe().PlayerStart = newStartPos
					c.Player.Position = newStartPos
					break
				}
			}
		}
		c.Location.Cafe().InEditorMode = false
		c.Location.SetRunning(true)
		log.Printf("SET RUNING TRUE")
		c.SendExtensionResponse("edi", "-1", "0",
			strconv.Itoa(status), strconv.Itoa(c.Player.Position[0]), strconv.Itoa(c.Player.Position[1]), "") // <- Objects
	case 1:
		c.SendExtensionResponse("edi", "-1", "0",
			strconv.Itoa(status), strconv.Itoa(c.Player.Position[0]), strconv.Itoa(c.Player.Position[1]), "") // <- Objects
		c.Location.Cafe().InEditorMode = true
		c.Location.SetRunning(false)
		log.Printf("SET RUNING FALSE")
	default:
		return nil
	}
	return nil
}
