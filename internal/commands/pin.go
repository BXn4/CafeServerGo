package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
)

// pin
func SendPing(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
  
  c.SendExtensionResponse("pin", "-1")

  return nil
}

/*
async def handle_pin(server: 'CafeServer', client: 'StreamWriter') -> None:
    logger.info("PIN PIN PIN PIN PIN")
    address = client.get_extra_info('peername')
    player = server.clients.get(address)

    if not player:
        client.close()

    while True:
        try:
            response = ExtensionResponse('pin', '-1')
            await server.send_response(client, response)
            await asyncio.sleep(60)
        except Exception as e:
            print(f"[ERROR] Cant send response, because the client is gone")
            print(f"[ERROR] An error occurred: {e}")
            break

    client.close()


*/
