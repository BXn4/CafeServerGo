package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
)

// ebf - C2S_EDITOR_BUY_FLOOR
func BuyFloor(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	// Dont allow players to modify the packet and sending us EBF while not in editor.
	if !c.Location.Cafe().InEditorMode {
		return nil
	}

	cafeSize := c.Location.Cafe().Size

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
	playerHave := c.Location.Cafe().FurnitureInventory[tileID]
	buyAmount, oldTiles := c.Location.Cafe().GetOldTiles(startX, startY, endX, endY, tileID)
	buyAmount = buyAmount - playerHave

	tileInfo, err := utils.GetTile(tileID)
	if err != nil {
		return err
	}

	// If the player not have enough cash
	if c.Player.Cash < tileInfo.Cash*buyAmount {
		c.SendExtensionResponse("ebf", "-1", "4",
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
		existingTileID := c.Location.Cafe().Tiles[tile[0]][tile[1]]

		fmt.Println("TILE ID: ", existingTileID)
		fmt.Println("AMOUNT: ", c.Location.Cafe().FurnitureInventory[existingTileID])

		// Add replaced tile to inventory
		if c.Location.Cafe().FurnitureInventory[existingTileID] != 0 {
			c.Location.Cafe().FurnitureInventory[existingTileID] += amount
		} else {
			c.Location.Cafe().FurnitureInventory[existingTileID] = amount
		}

		// Replace tile
		c.Location.Cafe().Tiles[tile[0]][tile[1]] = tileID
		oldTilesStr = append(oldTilesStr, fmt.Sprintf("%v+%v", existingTileID, amount))
	}

	c.SendExtensionResponse("ebf", "-1", "0",
		strconv.Itoa(startX),
		strconv.Itoa(startY),
		strconv.Itoa(endX),
		strconv.Itoa(endY),
		strconv.Itoa(tileID),
		strings.Join(oldTilesStr, "#"),
		strconv.Itoa(playerHave),
		strconv.Itoa(buyAmount),
	)

	return nil
}
