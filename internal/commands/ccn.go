package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

func Clean(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us CCH while in editor.
	if c.Location.Cafe().InEditorMode() {
		return nil
	}

	objX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	objY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	obj := c.Location.Cafe().GetObjectByPos(objX, objY)

	if obj == nil {
		return err
	}

	var status string
	if obj.IsStove() {
		if c.Player.Cash < 15 {
			status = "4"
		} else {
			status = "0"
			c.Player.Cash -= 15
			obj.SetDishID(-1)
		}
	} else {
		obj.SetDishID(-1)
		status = "0"
	}

	c.SendExtensionResponse(
		"ccn", "-1",
		status,
		req.Args[2],
		req.Args[3],
		utils.If(obj.IsStove(), "1", "0"),
	)

	return nil
}

/*
async def handle_ccn(server: 'CafeServer', client: 'StreamWriter', *params: str):
    address = client.get_extra_info('peername')
    player = server.clients[address]

    obj_x = params[1]
    obj_y = params[2]

    obj = player.cafe.get_object(int(obj_x), int(obj_y))

    if obj.type == ObjectType.STOVE:
        if player.cash < 15:
            status = '4'
        else:
            status = '0'
            player.cash -= 15
            obj.dish_id = -1
            server.db.update_player(player.id, cash=player.cash)
            server.db.update_cafe(player.id, objects=player.cafe.get_objects_as_json())
    else:
        status = '0'
        obj.dish_id = -1
        server.db.update_cafe(player.id, objects=player.cafe.get_objects_as_json())

    response = ExtensionResponse('ccn', '-1', status, obj_x, obj_y, str(int(obj.type == ObjectType.STOVE)))
    await server.send_response(client, response)
*/
