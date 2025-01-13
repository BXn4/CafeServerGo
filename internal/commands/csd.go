package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"strconv"
)

func StoveDeliver(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	stoveX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	stoveY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	counterX, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}

	counterY, err := strconv.Atoi(req.Args[5])
	if err != nil {
		return err
	}

	stove := c.Location.Cafe().GetObjectByPos(stoveX, stoveY)
	counter := c.Location.Cafe().GetObjectByPos(counterX, counterY)

	// Add to counter
	if counter.DishID == stove.DishID {
		counter.DishAmount += stove.DishAmount
	} else {
		counter.DishID = stove.DishID
		counter.DishAmount = stove.DishAmount
	}

	// Get dish
	dish, err := utils.GetDish(stove.DishID)
	if err != nil {
		return err
	}

	// Reset stove
	stove.DishID = -2 // Dirty
	stove.FancyIng = false
	stove.StartedAt = nil
	stove.FinishesAt = nil

	// Increase xp
	c.Player.XP += dish.XP

	// Increase mastery
	// TODO: Create masteryies

	// response = ExtensionResponse('csd', '-1', '0', stove_x, stove_y, counter_x, counter_y, str(player.id))
	c.SendExtensionResponse(
		"csd", "-1", "0",
		req.Args[2],
		req.Args[3],
		req.Args[4],
		req.Args[5],
		strconv.Itoa(c.Player.ID),
	)
	// [RECEIVED] %xt%CafeEx%csd%-1%1%5%3%5%
	// [SENT] %xt%csd%-1%0%1%5%3%5%1%

	return nil
}

/*


2025-01-13 15:14:36 - INFO - server - [RECEIVED] %xt%CafeEx%csd%-1%1%5%3%6%
2025-01-13 15:14:36 - INFO - server - [SENT] %xt%csd%-1%0%1%5%3%6%2%
*/

/*

async def handle_csd(server: 'CafeServer', client: 'StreamWriter', *params: str) -> None:
    address = client.get_extra_info('peername')
    player = server.clients[address]

    stove_x = params[1]
    stove_y = params[2]
    counter_x = params[3]
    counter_y = params[4]

    stove = player.cafe.get_object(int(stove_x), int(stove_y))
    counter = player.cafe.get_object(int(counter_x), int(counter_y))

    dish_id = stove.dish_id

    stove.dish_id = -2
    stove.fancy_ing = False
    stove.started_at = None
    stove.finishes_at = None

    //----

    amount = player.get_mastery_amount(dish_id)
    xp = player.get_mastery_xp(dish_id)

    if counter.dish_id == dish_id:
        counter.dish_amount += amount
    else:
        counter.dish_id = dish_id
        counter.dish_amount = amount

    player.xp += xp

    server.db.update_player(player.id, xp=player.xp)
    server.db.update_cafe(player.id, objects=player.cafe.get_objects_as_json())

    response = ExtensionResponse('csd', '-1', '0', stove_x, stove_y, counter_x, counter_y, str(player.id))
    await server.send_response(client, response)

*/
