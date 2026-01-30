package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"

	"github.com/charmbracelet/log"
)

func init() {
	RegisterCommand(requests.POLICY_FILE,
		CommandConfig{
			Name:       "PolicyFileResponse",
			Identifier: "",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		PolicyFileResponse,
		nil,
	)

	RegisterCommand(requests.VERSION_CHECK,
		CommandConfig{
			Name:       "VersionCheckResponse",
			Identifier: "",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		VersionCheckResponse,
		nil,
	)

	RegisterCommand(requests.AUTO_JOIN,
		CommandConfig{
			Name:       "AutoJoinResponse",
			Identifier: "",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		AutoJoinResponse,
		nil,
	)

	RegisterCommand(requests.ROUND_TRIP,
		CommandConfig{
			Name:       "RoundTripResponse",
			Identifier: "",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		RoundTripResponse,
		nil,
	)

	RegisterCommand(requests.DISCONNECT,
		CommandConfig{
			Name:       "DisconnectResponse",
			Identifier: "",
			MinArgs:    0,
			MaxArgs:    0,
		},
		nil,
		DisconnectResponse,
		nil,
	)
}

func PolicyFileResponse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *CommandConfig) error {
	c.SendSystemResponse(responses.POLICY_FILE)
	return nil
}

func VersionCheckResponse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *CommandConfig) error {
	c.SendSystemResponse(responses.VERSION_CHECK)
	return nil
}

func AutoJoinResponse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *CommandConfig) error {
	c.SendSystemResponse(responses.AUTO_JOIN)
	return nil
}

func RoundTripResponse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *CommandConfig) error {
	c.SendSystemResponse(responses.ROUND_TRIP)
	return nil
}

func DisconnectResponse(req *requests.Request, c *client.Client, gm *managers.GameManager, cm *CommandConfig) error {
	log.Infof("Logout request from client %d, username: %v", c.ClientID, c.Player.GetUsername())

	// Send logout response first to ensure client gets the message
	c.SendSystemResponse(responses.LOGOUT)
	c.Writer.Flush() // Ensure the message is sent immediately

	// Immediately disconnect and cleanup player state
	log.Infof("Calling DisconnectClient for client %d", c.ClientID)
	gm.DisconnectClient(c.ClientID)
	log.Infof("DisconnectClient completed for client %d", c.ClientID)

	// Return an error to trigger return from HandleClient, closing connection
	return fmt.Errorf("client logged out")
}
