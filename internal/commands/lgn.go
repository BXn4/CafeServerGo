package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

// lgn - Login
func Login(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
	name := req.Args[2]
	password := req.Args[3]

	// Check credentials
	statusCode, err := c.DB.Authenticate(name, password)
	if err != nil {
		return err
	}
	statusCodeStr := strconv.Itoa(statusCode)

	// TODO Check if already logged in log him/her out

	// Send login response (lgn)
	c.SendExtensionResponse("lgn", "1", statusCodeStr)

	// Send room list (rlu)
	RoomList(req, c, clientManager, cafeManager)

	// TODO: Daily login reward check

	// Send user info (gui)
	err = UserInfo(req, c, clientManager, cafeManager)
	if err != nil {
		return err
	}

	// Send balancing constants (sbc)
	SendBalancingConstant(req, c, clientManager, cafeManager)

	// Send mastery info (lmi)
	SendMasteryInfo(req, c, clientManager, cafeManager)

	// Send fridge info (ifr)
	SendFridgeInventory(req, c, clientManager, cafeManager)

	// TODO: Handle login bonus (lbu)

	// Send Ping (pin)
	SendPing(req, c, clientManager, cafeManager)

	return nil
}

/*
async def handle_lgn(server: 'CafeServer', client: 'StreamWriter', *params: str) -> None:
    name_or_mail = params[1]
    password = params[2]

    status_code, username = server.db.try_login(name_or_mail, password)

    if status_code != 0:
        response = ExtensionResponse('lgn', '1', str(status_code))
        await server.send_response(client, response)
    else:
        // Check if logged in
        check_player = [player for player in server.players if player.avatar.username == username]
        if check_player:
            player = check_player[0]
        else:
            player = server.db.get_player(username)
            server.players.append(player)

        if player.online:
            response = ExtensionResponse('lgn', '1', '15')
            await server.send_response(client, response)

            return



        address = client.get_extra_info('peername')
        server.clients[address] = player

        await server.send_response(client, SystemResponse('logout'))

        response = ExtensionResponse('lgn', '1', '0')
        await server.send_response(client, response)


        await handle_rlu(server, client, 'lgn', *params)

        // HANDLE daily login

        if not player.daily_login:
            current_time = datetime.now(timezone.utc).strftime("%y-%d-%m %H:%M")
            player.daily_login = current_time
            server.db.update_player(player.id, daily_login=player.daily_login)

        # RESETS WHEN 24H PAST
        # !! NEED TO CALL BEFORE GUI !!
        if player.check_daily_login():
            player.instant_cookings = 0
            player.played_wheel = False
            player.open_jobs = 5
            server.db.update_player(player.id, instant_cookings=0, played_wheel=0, open_jobs=5)


        await handle_gui(server, client, *params)

        await handle_sbc(server, client, *params)

        await handle_lmi(server, client, *params)

        await handle_ifr(server, client, *params)

        await handle_lbu(server, client)

        await handle_pin(server, client)
*/
