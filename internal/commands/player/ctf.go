/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

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

func TutorialComplete(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *commands.CommandConfig) error {
	println("TUTORIAL DONE")
	println(c.Player.GetXP())
	if !c.Player.GetIsTutorialCompleted() {
		println("DONE")

		c.Player.AddXP(5) // 5 XP reward for tutorial complete

		go agents.StartAgentCycles(c.Location)

	}
	return nil
}
