/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package gift

import (
	"cafego/internal/utils"
	"fmt"
	"strings"
	"time"
)

type Gift struct {
	ID       int
	Amount   int
	Sender   int
	Received time.Time
}

func (g *Gift) String() string {
	return fmt.Sprintf("%v+%v+%v+%v", g.ID, g.Amount, g.Sender, g.Received.Format("2006-01-02"))
}

func NewGift(id, amount, sender int, received time.Time) *Gift {
	return &Gift{
		ID:       id,
		Amount:   amount,
		Sender:   sender,
		Received: received,
	}
}

func NewGiftFromString(s string) (*Gift, error) {

	data := strings.Split(s, "+")

	items, err := utils.MultiAtoi(data[:len(data)-1]...)
	if err != nil {
		return nil, err
	}
	println("NewGiftFromString: ", s)
	id, amount, sender := items[0], items[1], items[2]

	received, err := time.Parse("2006-01-02", data[3])
	if err != nil {
		return nil, err
	}

	return &Gift{
		ID:       id,
		Amount:   amount,
		Sender:   sender,
		Received: received,
	}, nil
}
