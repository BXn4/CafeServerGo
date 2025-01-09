package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// jca - JoinCafe
func JoinCafe(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	// Get id of cafe to join
	id, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	// Adds cafe to manager (load it if not loaded)
	cafe := cafeManager.Add(id)

	// Send cafe joined
	c.SendExtensionResponse("jca", "-1", "0")

	// Leave cafe if already in one
  if c.Cafe != nil {
    c.Cafe.Leave(c.Player.ID)
  }

	// Join cafe
	cafe.Join(id, c.Conn)

	// Save location
	c.Cafe = cafe

	return nil
}

