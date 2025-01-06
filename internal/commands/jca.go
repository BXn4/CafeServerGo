package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
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

  // TODO: Handle if already in a cafe

  // Join cafe 
  cafe.Join(id, c.Conn) 

  // Save location
  c.Cafe = cafe 

  // Send cafe layout (sgc)
  SendCafe(req, c, clientManager, cafeManager)



  return nil
}


/*
        existing_player = [player for player in server.players if player.id == int(params[2])]
        if existing_player:
            wanted_player = existing_player[0]
        else:
            wanted_player = server.db.get_player_by_id(int(params[2]))
            server.players.append(wanted_player)

        response = wanted_player.cafe.to_response('sgc', '-1', '0')
        player.room = wanted_player.cafe
        player.pos = wanted_player.cafe.get_start_pos()

        await server.send_response(client, ExtensionResponse(*response))

*/



/*
async def handle_jca(server: 'CafeServer', client: 'StreamWriter', *params: str) -> None:
    wanted_player_id = int(params[2])
    address = client.get_extra_info('peername')
    player = server.clients[address]

    check_player = [player for player in server.players if player.id == wanted_player_id]
    if check_player:
        wanted_player = check_player[0]
    else:
        wanted_player = server.db.get_player_by_id(int(params[2]))
        server.players.append(wanted_player)

    response = ExtensionResponse('jca', '-1', '0')
    await server.send_response(client, response)


--------------------------\/\/\/\/\/\/


    await handle_sgc(server, client, 'jca', *params)

    await handle_ifr(server, client, *params)
    await handle_ein(server, client, *params)

    if wanted_player.id == player.id:
        await handle_asy(server, client)

    if not wanted_player.cafe.customer_cycle:
        print(f'Created a new customer cycle\nUsed this player\'s Cafe: {wanted_player.avatar.username}')
        wanted_player.cafe.customer_cycle = create_task(spawn_cycle(server, wanted_player.client))

        create_task(spawn_waiters(server, wanted_player.client))
    else:
        create_task(spawn_customers_to_client(server, client, wanted_player.cafe))
        create_task(spawn_waiters_to_client(server, client, wanted_player.cafe))
*/

