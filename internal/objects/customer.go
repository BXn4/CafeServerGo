package objects

import "sync"

type CustomerAction int

const (
	CUSTOMER_INSERT               CustomerAction = 0
	CUSTOMER_WALK_TO_CHAIR                       = 1
	CUSTOMER_SIT_DOWN                            = 2
	CUSTOMER_EAT                                 = 3
	CUSTOMER_LEAVE                               = 4
	CUSTOMER_FAST_FOOD                           = 8
	CUSTOMER_GOTO_VENDING_MACHINE                = 9
	CUSTOMER_LEAVE_COMPLETE                      = 41
)

type Customer struct {
	id             int
	avatar         Avatar
	pos            [2]int
	dish           int
	action         CustomerAction
	isThirsty      bool
	assignedWaiter int
	mutex          sync.Mutex
}

func NewCustomer(id int, avatar Avatar, pos [2]int, dish int, action CustomerAction, isThirsty bool, assignedWaiter int) *Customer {
	return &Customer{
		id:             id,
		avatar:         avatar,
		pos:            pos,
		dish:           dish,
		action:         action,
		isThirsty:      isThirsty,
		assignedWaiter: assignedWaiter,
	}
}

// --- GETTERS ---------------------
func (c *Customer) GetID() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.id
}

func (c *Customer) GetAvatar() *Avatar {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return &c.avatar
}

func (c *Customer) GetPos() [2]int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.pos
}

func (c *Customer) GetDish() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.dish
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

// --- SETTERS -----------------------

func (c *Customer) SetID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.id = id
}

func (c *Customer) SetAvatar(a Avatar) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.avatar = a
}

func (c *Customer) SetPos(pos [2]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pos = pos
}

func (c *Customer) SetDish(v int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.dish = v
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
