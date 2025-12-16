package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
)

type Command struct {
	Config    CommandConfig
	Validator func(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes)
	Handler   func(req *requests.Request, c *client.Client, gm *managers.GameManager) error
}

type CommandConfig struct {
	Identifier string
	Name       string
	MinArgs    int
	MaxArgs    int
	IsBool     bool
}

var Commands = map[requests.RequestKind]Command{}

func RegisterCommand(
	kind requests.RequestKind,
	commandConfig CommandConfig,
	commandValidator func(*requests.Request, *client.Client, *managers.GameManager, CommandConfig) (string, ErrorCodes),
	commandHandler func(*requests.Request, *client.Client, *managers.GameManager) error,
) {
	Commands[kind] = Command{
		Config:    commandConfig,
		Validator: commandValidator,
		Handler:   commandHandler,
	}
}

func ErrorHandler(req *requests.Request, c *client.Client, cm *CommandConfig, reason string, errc ErrorCodes) error {
	c.SendExtensionResponse(cm.Identifier, "-1", strconv.Itoa(int(errc)), strings.Join(req.Args[2:], "%"))
	return fmt.Errorf("Command %s failed with error code: %d.\n---> Reason: %s", cm.Name, errc, reason)
}

func HandleClient(c *client.Client, gm *managers.GameManager) {
	defer c.Disconnect()
	for req := range c.RequestQueue {
		if req == nil {
			return
		}

		if req.NeedsLogin() && c.Player == nil &&
			req.Kind != requests.C2S_LOGIN && req.Kind != requests.C2S_SPECIAL_EVENT {
			// While the player is not logged in, disconnects the client if the request is not for login.
			// The client sends SEE command after login. Maybe we can patch this in the game client to only send it, when the login was successful.
			return
		}

		// Handle requests
		err := HandleRequest(req, c, gm)
		if err != nil {
			log.Warnf(err.Error())
			continue
		}
	}

}

func HandleRequest(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	command, implemented := Commands[req.Kind]

	if !implemented {
		cm := CommandConfig{
			Name:       req.Args[0],
			Identifier: req.Args[0],
		}
		return ErrorHandler(req, c, &cm, "The command is not implemented", NOT_IMPLEMENTED)
	}

	// command error = int
	// all error codes what the cafe having in int
	if command.Validator != nil {
		reason, commandError := command.Validator(req, c, gm, command.Config)

		// If theres an ANY error, then dont run the handler.
		if commandError != NO_ERROR {
			return ErrorHandler(req, c, &command.Config, reason, commandError)
		}
	}

	return command.Handler(req, c, gm)
}
