package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
)

func TutorialComplete(c *client.Client, gm *managers.GameManager) error {
	if !c.Player.IsTutorialCompleted {
		c.Player.IsTutorialCompleted = true
		go agents.StartAgentCycles(c.Location)

	}
	return nil
}
