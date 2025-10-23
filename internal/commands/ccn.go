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
	if !c.Location.IsRunning() {
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

	obj := c.Location.Cafe().GetObjectByPosXY(objX, objY)

	if obj == nil {
		return err
	}

	if obj.IsStove() {
		if c.Player.GetCash() < 15 {
			c.SendExtensionResponse(
				"ccn", "-1",
				"4",
				req.Args[2],
				req.Args[3],
				utils.If(obj.IsStove(), "1", "0"),
			)

			return nil
		}

		c.Player.AddCash(-15)

		if obj.GetDishID() > 0 && obj.GetIsRotten() {
			c.Player.UpdateAchivementOvercookedFoods() // if the player cleans rotten food

			c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())
		}
		obj.SetDishID(-1)
	}

	if obj.IsCounter() {
		if obj.GetDishID() > 0 {
			obj.SetDishID(-1)
		}
	}

	c.SendExtensionResponse(
		"ccn", "-1",
		"0",
		req.Args[2],
		req.Args[3],
		utils.If(obj.IsStove(), "1", "0"),
	)

	c.DB.UpdateObjects(c.Location.Cafe().ID, c.Location.Cafe().Objects.StringForDB())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())

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
