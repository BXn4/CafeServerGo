package player

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

func init() {
	commands.RegisterCommand(requests.C2S_CAFE_TUTORIAL_FINISH,
		commands.CommandConfig{
			Name:       "TutorialFinish",
			Identifier: "ctf",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		TutorialComplete,
		nil,
	)
}

func TutorialComplete(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	if !c.Player.GetIsTutorialCompleted() {
		c.Player.SetIsTutorialCompleted(true)

		c.Player.AddXP(5) // 5 XP reward for tutorial complete

		go agents.StartAgentCycles(c.Location)

	}
	return nil
}
