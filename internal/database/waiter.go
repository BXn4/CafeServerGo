package database

import (
	"cafego/internal/objects"
	"cafego/internal/types/daos"
	"cafego/internal/utils"
	"strings"
)

func ConvertWaiterDAOToWaiter(dao *daos.WaiterDAO) (*objects.Waiter, error) {
	var waiter objects.Waiter

	// Fill simple waiter data
	waiter.ID = dao.ID
	waiter.Name = dao.Name
	waiter.Priority = dao.Priority

	// Parse avatar
	println("waiter avatar: ", dao.Avatar)
	data := strings.Split(dao.Avatar, "+")
	waiter.Avatar = *objects.NewAvatarFromString(data[2])
	waiter.Avatar.Name = data[0]
	waiter.Avatar.Gender = objects.AvatarGender(utils.If(data[1] == "2", 2, 1))
	waiter.Avatar.IsNPC = true

	return &waiter, nil
}
