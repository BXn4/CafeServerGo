package commands

import (
  "cafego/internal/types/requests"
  "cafego/internal/client"
  "cafego/internal/managers"
  "fmt"
  "strings"
  "strconv"
)

// ifr - SendFridgeInventory
func SendFridgeInventory(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {
 
  var fridge map[int]int

  if c.Cafe == nil {
    cafe, err := c.DB.GetCafeByPlayerID(c.Player.ID)
    if err != nil {
      return err
    }
    fridge = cafe.FridgeInventory
  }else{
    c.Cafe.Fridge()
  }

  fridgeCap := len(fridge)

  var fridgeArgs []string

  for k, v := range fridge {
    item := fmt.Sprintf("%v+%v", k, v)
    fridgeArgs = append(fridgeArgs, item)
  }
  
  c.SendExtensionResponse("ifr", "1", "0",
    strconv.Itoa(fridgeCap),
    strings.Join(fridgeArgs, "#"),
  )
  return nil
}

/*
async def handle_ifr(server: 'CafeServer', client: 'StreamWriter', *params: str) -> None:
    address = client.get_extra_info('peername')
    player = server.clients[address]

    fridge_capacity = player.cafe.fridge_capacity
    fridge_inventory = player.cafe.build_fridge_inventory()

    response = ExtensionResponse('ifr', '-1', '0', str(fridge_capacity), fridge_inventory)
    await server.send_response(client, response)


*/
