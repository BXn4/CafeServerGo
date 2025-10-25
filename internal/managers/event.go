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

func (gm *GameManager) SetDaysLeft(value int) {
	gm.gameEventDaysLeft = value
}

func (gm *GameManager) GetDaysLeft() int {
	return gm.gameEventDaysLeft
}

func (gm *GameManager) CheckForEvent(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	gm.Check()

	for range ticker.C {
		gm.Check()
	}
}

func (gm *GameManager) Check() {
	currentTime := time.Now().UTC()

	isEvent := utils.IsEvent(currentTime)

	switch isEvent {
	case true:
		if gm.GetEvent() == 0 {
			eventType := utils.GetEventType(currentTime)
			gm.SetEvent(eventType)
		}

		daysLeft := utils.GetDaysLeft(currentTime)
		gm.SetDaysLeft(daysLeft)
	default:
		gm.SetEvent(0)
		gm.SetDaysLeft(0)
	}
}
