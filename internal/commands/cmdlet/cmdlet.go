package cmdlet

import (
	"cafego/internal/client"
	"cafego/internal/managers"

	"fmt"
	"strings"
)

type AccessLevel int

const (
	User AccessLevel = iota
	Moderator
	Admin
)

type Cmdlet struct {
	level AccessLevel
	fn    func(*client.Client, *managers.GameManager, []string) error
}

func (cmd *Cmdlet) Run(c *client.Client, gm *managers.GameManager, args []string) error {
	if cmd.level > AccessLevel(c.Player.AccessLevel) {
		return fmt.Errorf("Access denied!")
	}
	return cmd.fn(c, gm, args)
}

var cmdlets = map[string]Cmdlet{
	"ach": {Admin, SetAchivement},
	"xp":  {Admin, SetXP},
}

func HandleCmdlets(c *client.Client, gm *managers.GameManager, message string) {

	args := strings.Split(message, " ")
	var err error

	cmd := args[0][1:]
	cmdlet, ok := cmdlets[cmd]
	if !ok {
		err = fmt.Errorf("Unknown command: %v", cmd)
	} else {
		err = cmdlet.Run(c, gm, args[1:])
	}

	if err != nil {
		c.SendExtensionResponse("cpu", "-1", "2", "0", err.Error())
	}
}
