package managers

import (
	"cafego/internal/utils"
	"time"
)

func (gm *GameManager) SetEvent(event int) {
	gm.gameEvent = event
}

func (gm *GameManager) GetEvent() int {
	return gm.gameEvent
}

func (gm *GameManager) CheckForEvent(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
		currentTime := time.Now().UTC()

		isEvent := utils.IsEvent(currentTime)
		if isEvent && gm.GetEvent() == 0 {
			eventType := utils.GetEventType(currentTime)
			if eventType != 0 {
				gm.SetEvent(eventType)
			}
		}

		daysLeft := utils.GetDaysLeft(currentTime)
		if daysLeft < 0 {
			gm.SetEvent(0)
		}
	}
}
