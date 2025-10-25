package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"fmt"
)

// see - S2C_SPECIAL_EVENT
func SendSpecialEvent(c *client.Client, gm *managers.GameManager) error {
	if gm.GetEvent() == 0 {
		return nil
	}

	c.SendExtensionResponse(
		"see",
		"-1",
		"0",
		fmt.Sprintf("%d#", gm.GetEvent()),
		fmt.Sprintf("%d+%d", gm.GetEvent(), gm.GetDaysLeft()))

	return nil
}
