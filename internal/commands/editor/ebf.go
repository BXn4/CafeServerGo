package editor

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	commands.RegisterCommand(requests.C2S_EDITOR_BUY_FLOOR,
		commands.CommandConfig{
			Name:        "BuyFloor",
			Identifier:  responses.S2C_EDITOR_BUY_FLOOR,
			Description: "Buying floor",
			Args:        "{startX} {startY} {endX} {endY} {tileID} {oldTilesSTR}#  {playerHave}  {buyAmount}",
			MinArgs:     10,
			MaxArgs:     10,
		},
		nil,
		nil,
		nil,
	)
}

// ebf - C2S_EDITOR_BUY_FLOOR
func BuyFloor(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {

	// Dont allow players to modify the packet and sending us EBF while not in editor.
	if c.Location.IsRunning() {
		return nil
	}

	cafeSize := c.Location.Cafe().GetSize()

	items, err := utils.MultiAtoi(req.Args[2:]...)
	if err != nil {
		return err
	}

	for _, item := range items {
		if item < 0 || item > cafeSize {
			return fmt.Errorf("Cant add floor, because the wrong pos")
		}
	}

	startX, startY, endX, endY, tileID := items[0], items[1], items[2], items[3], items[4]

	if startX > endX {
		startX, endX = endX, startX
	}
	if startY > endY {
		startY, endY = endY, startY
	}

	// If the player have some in their inventory, dont buy that amount
	playerHave := c.Location.Cafe().GetFurnitureInventory()[tileID]
	buyAmount, oldTiles := c.Location.Cafe().GetOldTiles(startX, startY, endX, endY, tileID)
	buyAmount = buyAmount - playerHave

	tileInfo, err := utils.GetTile(tileID)
	if err != nil {
		return err
	}

	// If the player not have enough cash
	if c.Player.GetCash() < tileInfo.Cash*buyAmount {
		c.SendExtensionResponse(cm.Identifier, "-1", "4",
			strconv.Itoa(startX),
			strconv.Itoa(startY),
			strconv.Itoa(endX),
			strconv.Itoa(endY),
			strconv.Itoa(tileID),
		)
		return nil
	}

	var oldTilesStr []string
	for tile, amount := range oldTiles {
		existingTileID := c.Location.Cafe().GetTiles()[tile[0]][tile[1]]

		// Add replaced tile to inventory
		c.Location.Cafe().AddFurnitures(existingTileID, amount)

		// Replace tile
		c.Location.Cafe().SetTile(tile[0], tile[1], tileID)
		oldTilesStr = append(oldTilesStr, fmt.Sprintf("%v+%v", existingTileID, amount))
	}

	c.SendExtensionResponse(cm.Identifier, "-1", "0",
		strconv.Itoa(startX),
		strconv.Itoa(startY),
		strconv.Itoa(endX),
		strconv.Itoa(endY),
		strconv.Itoa(tileID),
		strings.Join(oldTilesStr, "#"),
		strconv.Itoa(playerHave),
		strconv.Itoa(buyAmount),
	)

	c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	c.DB.UpdateObjects(c.Location.Cafe().GetID(), c.Location.Cafe().GetObjects().StringForDB())

	return nil
}
