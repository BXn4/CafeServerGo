/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package customer

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
)

type CustomerAction int

const (
	CUSTOMER_INSERT               CustomerAction = 0
	CUSTOMER_WALK_TO_CHAIR        CustomerAction = 1
	CUSTOMER_SIT_DOWN             CustomerAction = 2
	CUSTOMER_EAT                  CustomerAction = 3
	CUSTOMER_LEAVE                CustomerAction = 4
	CUSTOMER_FAST_FOOD            CustomerAction = 8
	CUSTOMER_GOTO_VENDING_MACHINE CustomerAction = 9
	CUSTOMER_LEAVE_COMPLETE       CustomerAction = 41
)

type Customer struct {
	id             int
	avatar         avatar.Avatar
	pos            simple.Position
	action         CustomerAction
	isThirsty      bool
	assignedWaiter int
	dishID         int
	mutex          sync.Mutex
}

func NewCustomer(id int, avatar avatar.Avatar, pos simple.Position, dish int, action CustomerAction, isThirsty bool, assignedWaiter int) *Customer {
	avatar.IsNPC = true
	return &Customer{
		id:             id,
		avatar:         avatar,
		pos:            pos,
		action:         action,
		isThirsty:      isThirsty,
		assignedWaiter: assignedWaiter,
	}
}

func NewRandomCustomer(id int, pos simple.Position) *Customer {
	return &Customer{
		id:             id,
		avatar:         avatar.NewRandomAvatar(),
		pos:            pos,
		action:         CUSTOMER_INSERT,
		isThirsty:      rand.Float64() <= 0.05, // 5 % to spawn a thirsty cusotmer
		dishID:         -1,
		assignedWaiter: -1,
	}
}

// nac - npc avatar string
func (c *Customer) SpawnString() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	args := []string{
		strconv.Itoa(c.id),
		"0", // NPC type (0: Customer)
		"0",
		strconv.Itoa(c.dishID),
		utils.If(c.isThirsty, "1", "0"),
		c.avatar.String(),
	}

	return strings.Join(args, "+")
}

// nav - npc action string
func (c *Customer) ActionString() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	args := []string{
		strconv.Itoa(c.id),
		strconv.Itoa(int(c.action)),
	}
	// If there is no action return
	if c.action == CUSTOMER_INSERT || c.action == CUSTOMER_LEAVE {
		return strings.Join(args, "+")
	}

	args = append(args, strconv.Itoa(c.pos.X))
	args = append(args, strconv.Itoa(c.pos.Y))

	return strings.Join(args, "+")
}

// Modify the customer action to the visitor
// If the customer is walking, dont send walking to the visitor, send sit down
func (c *Customer) ActionStringToSpawnBack(action CustomerAction) string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	args := []string{
		strconv.Itoa(c.id),
		strconv.Itoa(int(action)),
	}

	if action == CUSTOMER_INSERT || action == CUSTOMER_LEAVE {
		return strings.Join(args, "+")
	}

	args = append(args, strconv.Itoa(c.pos.X))
	args = append(args, strconv.Itoa(c.pos.Y))

	return strings.Join(args, "+")
}

// --- GETTERS ---------------------
func (c *Customer) GetID() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.id
}

func (c *Customer) GetAvatar() *avatar.Avatar {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return &c.avatar
}

func (c *Customer) GetPos() simple.Position {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.pos
}

func (c *Customer) GetAction() CustomerAction {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.action
}

func (c *Customer) IsThirsty() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.isThirsty
}

func (c *Customer) GetAssignedWaiter() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.assignedWaiter
}

func (c *Customer) GetDishID() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.dishID
}

// --- SETTERS -----------------------

func (c *Customer) SetID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id = id
}

func (c *Customer) SetAvatar(a avatar.Avatar) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.avatar = a
}

func (c *Customer) SetPos(pos simple.Position) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pos = pos
}

func (c *Customer) SetAction(v CustomerAction) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.action = v
}

func (c *Customer) SetIsThirsty(v bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.isThirsty = v
}

func (c *Customer) SetAssignedWaiter(v int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.assignedWaiter = v
}

func (c *Customer) SetDishID(dishID int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dishID = dishID
}
