package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
)

// vck - VersionCheck
func VersionCheck(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
  c.SendExtensionResponse("vck", "1", "0", "1603")
  return nil
}
