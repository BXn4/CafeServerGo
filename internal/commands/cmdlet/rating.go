package cmdlet

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"fmt"
	"strconv"
)

func SetRating(c *client.Client, gm *managers.GameManager, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("/rating <new rating>")
	}

	rating, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	c.Location.Cafe().SetRating(rating)

	return nil
}
