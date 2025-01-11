package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

// ccc - C2S_CAFE_COOK
func StartCooking(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	posX, err := strconv.Atoi(req.Args[2])
	if err != nil {
		fmt.Printf("Cant parse posX to int: %v", err)
		return err
	}
	posY, err := strconv.Atoi(req.Args[3])
	if err != nil {
		fmt.Printf("Cant parse posY to int: %v", err)
		return err
	}
	dishID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		fmt.Printf("Cant parse dishID to int: %v", err)
		return err
	}
	isPrepared, err := strconv.Atoi(req.Args[5])
	if err != nil {
		fmt.Printf("Cant parse isPrepared to int: %v", err)
		return err
	}
	usingFancy, err := strconv.Atoi(req.Args[6])
	if err != nil {
		fmt.Printf("Cant parse usingFancy to int: %v", err)
		return err
	}

	cookingTime := c.Player.GetDishMasteryDuration(dishID)

	c.Location.Broadcast("ccc", "-1", "0",
		strconv.Itoa(posX), strconv.Itoa(posY), strconv.Itoa(dishID),
		strconv.Itoa(isPrepared), strconv.Itoa(usingFancy), strconv.Itoa(int(cookingTime)))

	return nil
}
