package waiter

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/customer"
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

const (
	WAITER_INSERT Action = iota
	WAITER_MOVE_TO_COUNTER
	WAITER_FEED
	WAITER_CLEAN
)

type Priority int

const (
	CLEANING = 0
	BOTH     = 50
	SERVING  = 100
)

type Action int

const (
	INSERT          Action = 0
	MOVE_TO_COUNTER Action = 5
	FEED            Action = 6
	CLEAN           Action = 7
)

type Waiter struct {
	id              int
	priority        int
	avatar          avatar.Avatar
	pos             simple.Position
	dish            int
	action          Action
	currentCounter  *object.Object
	currentCustomer *customer.Customer
	isWorking       bool
	mutex           sync.RWMutex
}

func NewWaiter(id, priority int, a avatar.Avatar, isWorking bool) *Waiter {
	return &Waiter{
		id:        id,
		priority:  priority,
		avatar:    a,
		isWorking: isWorking,
	}
}

func GetStartingWaiter() WaiterList {
	girlNames := []string{"Jane", "Lucy", "Emma", "Stacey", "Becky"}
	boyNames := []string{"Oskar", "James", "Jeffrey", "Tom", "William"}

	wa := NewWaiter(0, 50, avatar.NewRandomAvatar(), false)

	if wa.avatar.Gender == avatar.Girl {
		wa.avatar.Name = girlNames[rand.Intn(len(girlNames))]
	} else {
		wa.avatar.Name = boyNames[rand.Intn(len(boyNames))]
	}

	return WaiterList{wa}
}

func (w *Waiter) StopWorking() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.isWorking = false
	if w.currentCustomer != nil {
		w.currentCustomer.SetAssignedWaiter(-1)
	}
}

func (w *Waiter) String() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return fmt.Sprintf("%v+%v+%v", w.avatar.Name, w.avatar.Apperance(), w.priority)
}

// nac - npc avatar string
func (w *Waiter) SpawnString() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	args := []string{
		strconv.Itoa(w.id),
		"1", // NPC type (1: Waiter)
		strconv.Itoa(w.priority),
		"-1", // DishID
		w.avatar.String(),
	}
	return strings.Join(args, "+")
}

// nav - npc action string
func (w *Waiter) ActionString() string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	args := []string{
		strconv.Itoa(w.id),
		strconv.Itoa(int(w.action)),
		strconv.Itoa(w.GetPos().X),
		strconv.Itoa(w.GetPos().Y),
	}
	return strings.Join(args, "+")
}

// nav - npc action string
func (w *Waiter) ActionStringToSpawnBack(action Action, pos simple.Position) string {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	args := []string{
		strconv.Itoa(w.id),
		strconv.Itoa(int(action)),
		strconv.Itoa(pos.X),
		strconv.Itoa(pos.Y),
	}
	return strings.Join(args, "+")
}

func NewWaiterFromString(s string) *Waiter {
	parts := strings.Split(s, "+")

	priority, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}

	av := avatar.NewAvatarFromString(parts[1])
	av.IsNPC = true
	av.Name = parts[0]

	return &Waiter{
		avatar:    *av,
		priority:  priority,
		isWorking: false,
	}
}

// Getters
func (w *Waiter) GetID() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.id
}

func (w *Waiter) GetPriority() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.priority
}

func (w *Waiter) GetAvatar() avatar.Avatar {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.avatar
}

func (w *Waiter) GetPos() simple.Position {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.pos
}

func (w *Waiter) GetDish() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.dish
}

func (w *Waiter) GetAction() Action {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.action
}

func (w *Waiter) GetCurrentCounter() *object.Object {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.currentCounter
}

func (w *Waiter) GetCurrentCustomer() *customer.Customer {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.currentCustomer
}

func (w *Waiter) IsWorking() bool {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.isWorking
}

// Setters
func (w *Waiter) SetID(id int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.id = id
}

func (w *Waiter) SetPriority(priority int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.priority = priority
}

func (w *Waiter) SetAvatar(avatar avatar.Avatar) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.avatar = avatar
}

func (w *Waiter) SetPos(pos simple.Position) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.pos = pos
}

func (w *Waiter) SetDish(dish int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.dish = dish
}

func (w *Waiter) SetAction(action Action) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.action = action
}

func (w *Waiter) SetCurrentCounter(counter *object.Object) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.currentCounter = counter
}

func (w *Waiter) SetCurrentCustomer(customer *customer.Customer) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.currentCustomer = customer
}

func (w *Waiter) SetIsWorking(isWorking bool) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.isWorking = isWorking
}
