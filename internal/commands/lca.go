package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"fmt"
	"math/rand"
)

func CreateAvatar(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	guestName := fmt.Sprintf("Guest_%v'", rand.Intn(89999999)+10000000)

	avatar := objects.NewAvatarFromString(req.Args[2])
	if avatar == nil {
		return fmt.Errorf("Cant parse avatar from string!")
	}

	c.Player = &objects.Player{
		Avatar: *avatar,
	}

	c.SendExtensionResponse("lca", "-1", "0", guestName, "1")
	return nil
}
