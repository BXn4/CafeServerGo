package admin

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
	"strings"
)

func HandleAdminCommands(req *requests.Request, c *client.Client, gm *managers.GameManager, message string) {
	// TODO: check if admin

	args := strings.Split(message, " ")
	var err error
	switch args[0][1:] {
	case "ach":
		err = setAchivement(c.Player, args)
	}

	if err != nil {
		c.SendExtensionResponse("cpu", "-1", "2", "0", err.Error())
	}
}

func setAchivement(p *objects.Player, args []string) error {
	if len(args) != 3 {
		// TODO: Add help
		return fmt.Errorf("/ach <achivement id> <proggression>")
	}

	achivementID, err := strconv.Atoi(args[1])
	if err != nil {
		return err
	}

	proggression, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}

	if _, ok := p.Achievement[achivementID]; ok {
		p.Achievement[achivementID] = proggression
	}

	return nil
}
