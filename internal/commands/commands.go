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
	Handler   func(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) error
	DBSaver   func(c *client.Client) error
}

type CommandConfig struct {
	Name        string // Command name to us to easier to known which command is that (CafeWalk)
	Identifier  string // 3 letter identifier to the command (CWA)
	Description string // Whats the command will do.
	Args        string // What args needed to be send back.

	MinArgs int  // if the args less than X (needed) len, then dont allow command to run. Probaly will missing some values.
	MaxArgs int  // if the args more than X (needed) len, then dont allow command to run.
	IsBool  bool // 1 or 0

	PermissionLevel int // For the Staff / Admin / Moderator commands
	FeatureLevel    int // min level to use the command/feature. if its not declared, then its allow it.

	Category string // Categories. Like: Editor, Cafe, Player etc.
}

var Commands = map[requests.RequestKind]Command{}

func RegisterCommand(
	kind requests.RequestKind,
	commandConfig CommandConfig,
	commandValidator func(*requests.Request, *client.Client, *managers.GameManager, CommandConfig) (string, ErrorCodes),
	commandHandler func(*requests.Request, *client.Client, *managers.GameManager, CommandConfig) error,
	commandDBSaver func(*client.Client) error,
) {
	Commands[kind] = Command{
		Config:    commandConfig,
		Validator: commandValidator,
		Handler:   commandHandler,
		DBSaver:   commandDBSaver,
	}
}

func ErrorHandler(req *requests.Request, c *client.Client, cm *CommandConfig, reason string, errc ErrorCodes) error {
	c.SendExtensionResponse(cm.Identifier, "-1", strconv.Itoa(int(errc)), strings.Join(req.Args[2:], "%"))
	return fmt.Errorf("Command %s failed with error code: %d.\n---> Reason: %s", cm.Name, errc, reason)
}

func HandleClient(c *client.Client, gm *managers.GameManager) {
	for req := range c.RequestQueue {
		if req == nil {
			return
		}

		if req.NeedsLogin() && c.Player == nil &&
			req.Kind != requests.C2S_LOGIN &&
			req.Kind != requests.C2S_SPECIAL_EVENT &&
			req.Kind <= requests.VERSION_CHECK {
			// While the player is not logged in, disconnects the client if the request is not for login.
			// The client sends SEE command after login. Maybe we can patch this in the game client to only send it, when the login was successful.
			return
		}

		// Handle requests
		err := HandleRequest(req, c, gm)
		if err != nil {
			log.Errorf("Error during request handling: %v", err.Error())
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

	// log.Debugf("Handling command: %s", command.Config.Name)

	// command error = int
	// all error codes what the cafe having in int
	if command.Validator != nil {
		reason, commandError := command.Validator(req, c, gm, command.Config)

		// If theres an ANY error, then dont run the handler.
		if commandError != NO_ERROR {
			return ErrorHandler(req, c, &command.Config, reason, commandError)
		}
	}

	err := command.Handler(req, c, gm, command.Config)
	if err != nil {
		return fmt.Errorf("Error during command handling: %w", err)
	}

	if command.DBSaver != nil {
		err = command.DBSaver(c)
		if err != nil {
			return fmt.Errorf("Error during command DB saving: %w", err)
		}
	}

	return nil
}
