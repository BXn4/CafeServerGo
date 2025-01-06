package database

import (
	"cafego/internal/objects"
	"strconv"
	"strings"
)

type WaiterDAO struct {
	ID       int    `json:"waiter_id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Avatar   string `json:"avatar"`
}

func ConvertWaiterDAOToWaiter(dao *WaiterDAO) (*objects.Waiter, error) {
	var waiter objects.Waiter

	// Fill simple waiter data
	waiter.ID = dao.ID
	waiter.Name = dao.Name
	waiter.Priority = dao.Priority

	// Parse avatar
	waiter.Avatar = objects.Avatar{}
	data := strings.Split(dao.Avatar, "+")
	waiter.Avatar.Name = data[0]

	apperances := strings.Split(data[2], "#")
	for _, apperance := range apperances {

		// Parse
		values := strings.Split(apperance, "$")
		color, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, err
		}

		id := values[0][:len(values[0])-1]

		// Set values
		if values[0] == "1001" {
			waiter.Avatar.Gender = objects.Girl
			waiter.Avatar.TopColor = color
		} else if values[0] == "1002" {
			waiter.Avatar.Gender = objects.Boy
			waiter.Avatar.TopColor = color
		} else if id == "102" {
			waiter.Avatar.SkinColor = color
		} else if id == "104" {
			waiter.Avatar.HairColor = color
		} else if id == "105" {
			waiter.Avatar.LegsColor = color
		}
	}
	return &waiter, nil
}
