package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/avatar"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"fmt"
	"math/rand"

	"github.com/charmbracelet/log"
)

func CreateAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	guestName := fmt.Sprintf("Guest_%v", rand.Intn(89999999)+10000000)
	// userName+gender+Avatar+smartfoxClient.connectionTime+smartfoxClient.roundTripTime
	// user1   +2     +1052$12#1062$0#1042$9#1082$+0#1002$5#1022$5%1%60%1%1759766473155513873%
	log.Debugf("LCA Register avatar received: %s", req.Args[2])

	avatar := avatar.NewAvatarFromString(req.Args[2])
	avatar.Name = guestName

	if avatar == nil {
		c.SendExtensionResponse("lca", "-1", "-1", "-1", "-1")
		return fmt.Errorf("Cant parse avatar from string!")
	}

	if !avatar.IsValid() {
		c.SendExtensionResponse("lca", "-1", "-1", "-1", "-1")
		return fmt.Errorf("Invalid avatar structure!")
	}

	c.Player = &player.Player{
		Avatar: *avatar,
	}

	c.SendExtensionResponse("lca", "-1", "0", guestName, "1")
	return nil
}
