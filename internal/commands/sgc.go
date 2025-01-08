package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// sgc - Send cafe
func SendCafe(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	//args := []string{"sgc", "-1", "0"}
	//args = append(args, c.Cafe.AsResponse()...)

	//c.SendExtensionResponse(args...)

	//TODO: jul

	return nil
}

/*
async def handle_sgc(server: 'CafeServer', client: 'StreamWriter', cmd: str, *params: str) -> None:
    address = client.get_extra_info('peername')
    player = server.clients[address]

    if player.room:
        player.room.clients.remove(client)
        for user in player.room.clients:
            if user != client:
                await handle_juq(server, user, client)
                await handle_juq(server, client, user)


--------------\/\/\/\/\/\/\/\/\/

    -------------- JOIN ROOM -----------------------------------
    if cmd == 'jca':
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

    -------------- JOIN MARKETPLACE -----------------------------------
    elif cmd == 'mjm':
        response = server.marketplace.to_response('sgc', '-1', '0')
        player.room = server.marketplace
        player.pos = server.marketplace.player_start

        await server.send_response(client, ExtensionResponse(*response))

    --------------
    player.room.clients.append(client)

    # NOTE: *player.room.clients also includes the client that just joined the room
    await handle_jul(server, client, *player.room.clients)

    for user in player.room.clients:
        if user != client:
            await handle_juj(server, user, client)
*/
