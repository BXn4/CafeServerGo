package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"fmt"
)

// see - S2C_SPECIAL_EVENT
func SendSpecialEvent(c *client.Client, gm *managers.GameManager) error {
	if event.GetEvent() == 0 || c.Player == nil {
		return nil
	}

	c.SendExtensionResponse(
		"see",
		"-1",
		"0",
		fmt.Sprintf("%d#", event.GetEvent()),
		fmt.Sprintf("%d+%d", event.GetEvent(), event.GetDaysLeft()))

	return nil
}
