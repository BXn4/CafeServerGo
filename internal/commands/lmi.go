package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
)

// lmi - SendMasteryInfo
func SendMasteryInfo(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

  println("PLAYER MASTERY: ", c.Player.Mastery)

  c.SendExtensionResponse("lmi", "-1", "0", c.Player.Mastery)

  return nil
}


/*

async def handle_lmi(server: 'CafeServer', client: 'StreamWriter', *params: str) -> None:
    address = client.get_extra_info('peername')
    player = server.clients[address]

    response = ExtensionResponse('lmi', '-1', '0', player.build_mastery())
    await server.send_response(client, response)


*/
