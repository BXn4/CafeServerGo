package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// edit - C2S_EDITOR_MODE
func EditorMode(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// Parse status
	status, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	// Do not allow to enter editor mode multiple times
	if !c.Location.IsRunning() && status != 0 {
		return nil
	}

	switch status {
	case 0:
		c.Location.SetRunning(true)
		c.Location.ClearReservedObjects()
		c.Player.Position = c.Location.Cafe().GetPlayerStart()
	case 1:
		c.Location.SetRunning(false)
	}

	c.SendExtensionResponse("edi", "-1", "0",
		strconv.Itoa(status),
		strconv.Itoa(c.Player.Position.X),
		strconv.Itoa(c.Player.Position.Y),
		"", // <- Objects
	)

	return nil
}
