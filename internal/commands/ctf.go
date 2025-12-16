package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

func init() {
	RegisterCommand(requests.C2S_CAFE_TUTORIAL_FINISH,
		CommandConfig{
			Name:       "TutorialFinish",
			Identifier: "ctf",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		TutorialComplete,
	)
}

func TutorialComplete(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	if !c.Player.IsTutorialCompleted {
		c.Player.IsTutorialCompleted = true

		c.Player.AddXP(5) // 5 XP reward for tutorial complete

		SendSpecialEvent(req, c, gm)

		go agents.StartAgentCycles(c.Location)

	}
	return nil
}
