package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
)

// rlu - RoomList
func RoomList(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
  
  var roomID string
  if req.Args[0] == "lgn" {
    roomID = "1"
  } else {
    roomID = "-1"
  }

  c.SendExtensionResponse("rlu", roomID, "1", "1", "20", "2", "Lobby")
  return nil
}
/*
async def handle_rlu(server: 'CafeServer', client: 'StreamWriter', cmd: str, *params: str) -> None:

    if logged in send 1
    if cmd == 'lgn':
        room_id = '1'

    if not logged in send -1
    else:
        room_id = '-1'

    response = ExtensionResponse('rlu', room_id, '1', '1', '20', '2', 'Lobby')
    await server.send_response(client, response)
*/
