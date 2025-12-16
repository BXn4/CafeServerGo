package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
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
	)
}

func PolicyFileResponse(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendSystemResponse(responses.POLICY_FILE)
	return nil
}

func VersionCheckResponse(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendSystemResponse(responses.VERSION_CHECK)
	return nil
}

func AutoJoinResponse(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendSystemResponse(responses.AUTO_JOIN)
	return nil
}

func RoundTripResponse(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendSystemResponse(responses.ROUND_TRIP)
	return nil
}

func DisconnectResponse(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	c.SendSystemResponse(responses.LOGOUT)
	return nil
}
