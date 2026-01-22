package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/models/event"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
)

func init() {
	commands.RegisterCommand(requests.C2S_SPECIAL_EVENT,
		commands.CommandConfig{
			Name:        "SpecialEvent",
			Identifier:  responses.S2C_SPECIAL_EVENT,
			Description: "Tell the client if is there any event",
			Args:        "{0} {eventID}# {eventID}+{eventDaysLeft}",
			MinArgs:     0,
			MaxArgs:     0,
		},
		nil,
		SendSpecialEvent,
		nil,
	)
}

// see - S2C_SPECIAL_EVENT
func SendSpecialEvent(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	if event.GetEvent() == 0 || c.Player == nil {
		return nil
	}

	c.SendExtensionResponse(
		cm.Identifier,
		"-1",
		"0",
		fmt.Sprintf("%d#", event.GetEvent()),
		fmt.Sprintf("%d+%d", event.GetEvent(), event.GetDaysLeft()))

	return nil
}
