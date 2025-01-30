package database

import (
	"cafego/internal/objects"
	"strconv"
	"strings"
)

func NewWaiterFromString(s string) (*objects.Waiter, error) {
	var waiter objects.Waiter

	data := strings.Split(s, "+")
	priority, err := strconv.Atoi(data[2])
	if err != nil {
		return nil, err
	}

	// Fill simple waiter data
	waiter.ID = -1
	waiter.Name = data[0]
	waiter.Priority = priority

	// Parse avatar
	waiter.Avatar = *objects.NewAvatarFromString(data[1])
	waiter.Avatar.IsNPC = true

	return &waiter, nil
}
