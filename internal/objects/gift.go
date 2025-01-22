package objects

import (
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Gift struct {
	ID       int
	Amount   int
	Sender   int
	Received time.Time
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

func ParseGifts(s string) ([]*Gift, error) {
	if !strings.Contains(s, "#") {
		return []*Gift{}, nil
	}
	giftsStr := strings.Split(s, "#")
	gifts := []*Gift{}
	for _, giftStr := range giftsStr {
		gift, err := NewGiftFromString(giftStr)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}
	return gifts, nil
}

func (g *Gift) String() string {
	return fmt.Sprintf("%v+%v+%v+%v", g.ID, g.Amount, g.Sender, g.Received.Format("2006-01-02"))
}

func (p *Player) AddGift(id, amount, sender int) {
	gift := NewGift(id, amount, sender, time.Now().UTC())
	p.Gifts = append(p.Gifts, gift)
	p.NewGifts = len(p.Gifts)
}

func (p *Player) RemoveGift(index int) {
	p.Gifts = append(p.Gifts[:index], p.Gifts[index+1:]...)
	p.NewGifts = len(p.Gifts)
}

func BuildGifts(gifts []*Gift) string {

	giftsStr := []string{}
	for _, gift := range gifts {
		giftsStr = append(giftsStr, gift.String())
	}

	return strings.Join(giftsStr, "#")
}

func BuildGiftsWithIndex(gifts []*Gift) string {
	giftsStr := []string{}
	for i, gift := range gifts {
		iStr := strconv.Itoa(i)
		giftsStr = append(giftsStr, iStr+"+"+gift.String())
	}

	return strings.Join(giftsStr, "#")
}
