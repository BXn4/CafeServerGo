package daos

import (
	"strconv"
	"strings"
)

type WaiterDAO struct {
	ID       int    `json:"waiter_id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Avatar   string `json:"avatar"`
}

func NewWaiterDAOFromString(s string) (*WaiterDAO, error) {
	data := strings.Split(s, "+")
	id, err := strconv.Atoi(data[0])
	if err != nil {
		return nil, err
	}
	priority, err := strconv.Atoi(data[3])
	if err != nil {
		return nil, err
	}
	return &WaiterDAO{
		ID:       id,
		Name:     data[1],
		Priority: priority,
		Avatar:   data[2],
	}, nil
}
