package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/log"
)

func init() {
	RegisterCommand(requests.C2S_CREATE_AVATAR,
		CommandConfig{
			Name:       "CreateAvatar",
			Identifier: responses.S2C_CREATE_AVATAR,
			MinArgs:    7,
			MaxArgs:    7,
		},
		CreateAvatarValidator,
		CreateAvatar,
	)
}

func CreateAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	guestName := fmt.Sprintf("Guest_%v", rand.Intn(89999999)+10000000)
	// userName+gender+Avatar+smartfoxClient.connectionTime+smartfoxClient.roundTripTime
	// user1   +2     +1052$12#1062$0#1042$9#1082$+0#1002$5#1022$5%1%60%1%1759766473155513873%
	log.Debugf("LCA Register avatar received: %s", req.Args[2])

	avatar := avatar.NewAvatarFromString(req.Args[2])
	avatar.Name = guestName

	c.Player = &player.Player{
		Avatar: *avatar,
	}

	c.SendExtensionResponse("lca", "-1", "0", guestName, "1")
	return nil
}

func CreateAvatarValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	guestName := fmt.Sprintf("Guest_%v", rand.Intn(89999999)+10000000)
	avatar := avatar.NewAvatarFromString(req.Args[2])
	avatar.Name = guestName

	if avatar == nil {
		return "Cant parse avatar from the string, avatar is invalid!", NOT_DECLARED
	}

	if !avatar.IsValid() {
		return "Cant parse avatar from the string, avatar is invalid!", NOT_DECLARED
	}

	return "Command ran without any errors.", NO_ERROR
}
