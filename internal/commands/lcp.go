package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"

	"github.com/charmbracelet/log"
)

// lcp - ChangePassword
func ChangePassword(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	oldPasswd, newPasswd := req.Args[2], req.Args[3]

	status, err := c.DB.ChangePassword(c.ID(), oldPasswd, newPasswd)
	if err != nil {
		log.Warnf("Error changing password: %v", err)
	}

	c.SendExtensionResponse("lcp", "1", strconv.Itoa(status))
	c.Disconnect()
	return nil
}
