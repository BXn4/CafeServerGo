package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
	"strings"
)

// ifr - SendFridgeInventory
func SendFridgeInventory(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	var fridge map[int]int
	var fridgeCap int

	if c.Location == nil {
		cafe, err := c.DB.GetCafeByPlayerID(c.Player.ID)
		if err != nil {
			return err
		}
		fridge = cafe.FridgeInventory
		fridgeCap = cafe.GetFridgeMaxCapacity()
	} else {
		fridgeCap = c.Location.Cafe().GetFridgeMaxCapacity()
	}

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
