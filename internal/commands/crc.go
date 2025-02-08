package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// crc - Recook
func Recook(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// NOTE: This gives a chance to recover spoiled food for 1 gold
	// TODO
	//
	for _, arg := range req.Args {
		print(arg, " ")
	}
	println("")
	//c.SendExtensionResponse("crc", "1", "0", "1603")
	return nil
}
