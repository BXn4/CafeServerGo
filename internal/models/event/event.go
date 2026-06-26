/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package event

import (
	"cafego/internal/utils"
	"time"
)

var currentEvent = 0
var daysLeft = 0

func setEvent(event int) {
	currentEvent = event
}

func GetEvent() int {
	return currentEvent
}

func setDaysLeft(value int) {
	daysLeft = value
}

func GetDaysLeft() int {
	return daysLeft
}

func CheckForEvent(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	Check()

	for range ticker.C {
		Check()
	}
}

func Check() {
	currentTime := time.Now().UTC()

	isEvent := utils.IsEvent(currentTime)

	switch isEvent {
	case true:
		if GetEvent() == 0 {
			eventType := utils.GetEventType(currentTime)
			setEvent(eventType)
		}

		daysLeft := utils.GetDaysLeft(currentTime)
		setDaysLeft(daysLeft)
	default:
		setEvent(0)
		setDaysLeft(0)
	}
}
