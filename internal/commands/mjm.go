package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

//mjm - JoinMarketplace
func JoinMarketplace(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

  // Gets cafe location
	cafe := cafeManager.Add(-1)

	// Send cafe joined
	c.SendExtensionResponse("mjm", "-1", "0")

	// Leave current cafe if there is one
  if c.Cafe != nil {
    c.Cafe.Leave(c.Player.ID)
  }

	// Join cafe
	cafe.Join(c.Player.ID, c.Conn)

	// Save location
	c.Cafe = cafe

	return nil
}

