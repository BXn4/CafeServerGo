package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/utils"
	"fmt"
	"time"
)

// see - S2C_SPECIAL_EVENT
func SendSpecialEvent(c *client.Client, gm *managers.GameManager) error {

	currentTime := time.Now().UTC()

	isEvent := utils.IsEvent(currentTime)
	if isEvent && gm.GetEvent() == 0 {
		eventType := utils.GetEventType(currentTime)

		if eventType != 0 {
			gm.SetEvent(eventType)
		}

		daysLeft := utils.GetDaysLeft(currentTime)
		if daysLeft >= 0 {
			c.SendExtensionResponse(
				"see",
				"-1",
				"0",
				fmt.Sprintf("%d#", eventType),
				fmt.Sprintf("%d+%d", eventType, daysLeft),
			)
		} else {
			gm.SetEvent(0)
		}

	}

	return nil
}
