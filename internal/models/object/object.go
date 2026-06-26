/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package object

import (
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type Object struct {
	kind     CafeObjectKind
	pos      simple.Position
	rotation CafeObjectRotation

	dishID     int
	dishStatus int
	dishAmount int
	occupied   bool

	fancyIng   bool
	startedAt  *time.Time
	finishesAt *time.Time

	mutex sync.Mutex
}

func NewObject(posX int, posY int, objID int, objRotation int) (*Object, error) {
	cafeObj := Object{
		pos:      simple.NewPosition(posX, posY),
		kind:     CafeObjectKind(objID),
		rotation: CafeObjectRotation(objRotation),
	}
	return &cafeObj, nil
}

func NewObjectFromString(s string) (*Object, error) {
	parts := strings.Split(s, "+")
	items, err := utils.MultiAtoi(parts...)
	if err != nil || len(items) < 4 {
		return nil, fmt.Errorf("Invalid object string format: %v", err)
	}

	obj := &Object{
		pos:      simple.NewPosition(items[0], items[1]),
		kind:     CafeObjectKind(items[2]),
		rotation: CafeObjectRotation(items[3]),
		dishID:   -1,
	}

	if len(items) < 5 {
		return obj, nil
	}

	if items[4] <= 0 {
		obj.dishID = items[4]
		return obj, nil
	}

	if obj.isStove() {
		obj.dishID = items[4]
		obj.fancyIng = utils.If(items[5] == 1, true, false)
		if items[6] != -1 {
			startedAt, err := time.Parse(time.RFC3339, parts[6])
			finishesAt, err := time.Parse(time.RFC3339, parts[7])
			if err != nil {
				return nil, fmt.Errorf("Invalid object string format for stoves: %v", err)
			}
			obj.startedAt = &startedAt
			obj.finishesAt = &finishesAt
		}

	} else if obj.isCounter() || obj.isVending() {
		obj.dishID = items[4]
		obj.dishAmount = items[5]
	} else if obj.IsChair() {
		obj.dishID = items[4]
		obj.dishStatus = items[5]
		obj.occupied = false
	}

	return obj, nil
}

func (c *Object) IsWall() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 101 <= c.kind && c.kind <= 135
}

func (c *Object) IsWallObject() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 901 <= c.kind && c.kind <= 934
}

func (c *Object) IsDoor() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 201 <= c.kind && c.kind <= 207
}

func (c *Object) IsStove() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 252 <= c.kind && c.kind <= 259
}

func (c *Object) isStove() bool {
	return 252 <= c.kind && c.kind <= 259
}

func (c *Object) IsCounter() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 301 <= c.kind && c.kind <= 312
}

func (c *Object) isCounter() bool {
	return 301 <= c.kind && c.kind <= 312
}

func (c *Object) IsChair() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 601 <= c.kind && c.kind <= 624
}

func (c *Object) isChair() bool {
	return 601 <= c.kind && c.kind <= 624
}

func (c *Object) IsTable() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 401 <= c.kind && c.kind <= 423
}

func (c *Object) IsVending() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.kind == 1701
}

func (c *Object) isVending() bool {
	return c.kind == 1701
}

func (c *Object) IsFridge() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 351 <= c.kind && c.kind <= 358
}

func (c *Object) String() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	args := []string{
		strconv.Itoa(c.pos.X),
		strconv.Itoa(c.pos.Y),
		strconv.Itoa(int(c.kind)),
		strconv.Itoa(int(c.rotation)),
	}

	if c.isStove() {
		args = append(args, strconv.Itoa(c.dishID))
		if c.dishID > 0 {
			fancyIngStr := utils.If(c.fancyIng, "1", "0")
			args = append(args, fancyIngStr)
			if c.startedAt != nil {

				args = append(args,
					strconv.Itoa(c.GetPassedTime()),
					strconv.Itoa(c.GetRemaingTime()),
				)

			} else {
				args = append(args, "-1", "-1")
			}
		}
	} else if c.isCounter() || c.isVending() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishAmount))
	} else if c.isChair() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishStatus))
	}

	return strings.Join(args, "+")
}

func (c *Object) StringForDB() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	args := []string{
		strconv.Itoa(c.pos.X),
		strconv.Itoa(c.pos.Y),
		strconv.Itoa(int(c.kind)),
		strconv.Itoa(int(c.rotation)),
	}

	if c.isStove() {
		args = append(args, strconv.Itoa(c.dishID))
		if c.dishID > 0 {

			fancyIngStr := utils.If(c.fancyIng, "1", "0")
			args = append(args, fancyIngStr)

			if c.startedAt != nil {
				args = append(args,
					c.startedAt.Format(time.RFC3339),
					c.finishesAt.Format(time.RFC3339),
				)
			} else {
				args = append(args, "-1", "-1")
			}
		}
	} else if c.isCounter() || c.isVending() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishAmount))
	} else if c.isChair() {
		args = append(args, "-1", "0") // Not needed to save to DB
	}

	return strings.Join(args, "+")
}

func (c *Object) GetNormalizedRotation() [2]int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.rotation == Up {
		return [2]int{1, 0}
	} else if c.rotation == Down {
		return [2]int{-1, 0}
	} else if c.rotation == Left {
		return [2]int{0, -1}
	} else /* Right */ {
		return [2]int{0, 1}
	}
}

// --- GETTERS -------------------------------------

func (c *Object) GetKind() CafeObjectKind {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.kind
}

func (o *Object) GetPos() simple.Position {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	return o.pos
}

func (c *Object) GetRotation() CafeObjectRotation {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.rotation
}

func (c *Object) GetDishID() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishID
}

func (c *Object) GetDishStatus() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishStatus
}

func (c *Object) GetDishAmount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishAmount
}

func (c *Object) GetFancyIng() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.fancyIng
}

func (c *Object) GetStartedAt() *time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.startedAt
}

func (c *Object) GetPassedTime() int {
	return int(time.Now().UTC().Sub(*c.startedAt).Seconds())
}

func (c *Object) GetRemaingTime() int {
	return int(c.finishesAt.Sub(time.Now().UTC()).Seconds())
}

func (c *Object) GetFinishesAt() *time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.finishesAt
}

func (c *Object) GetOccupied() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.occupied
}

func (c *Object) GetIsRotten() bool {
	// Rotten Duration: return Math.max(CafeConstants.MIN_DISH_READY_TIME,this._baseDuration)
	// if not yet rotten: currentDish.rottenDuration * 60 + this.stoveVO.timeLeft > 0
	dishInfo, err := utils.GetDish(c.GetDishID())
	if err != nil {
		log.Printf("Invalid dish id: %s", err)
		return true
	}

	rottenDuration := math.Max(60, float64(dishInfo.Duration))
	if rottenDuration*60+float64(c.GetRemaingTime()) < 0 {
		// log.Debugf("The dish is rotten!")
		return true
	}

	// log.Debugf("The dish is not rotten!")

	return false
}

// --- SETTERS -------------------------------------

func (c *Object) SetKind(k CafeObjectKind) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.kind = k
}

func (c *Object) SetPos(pos simple.Position) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pos = pos
}

func (c *Object) SetPosXY(x, y int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pos.X = x
	c.pos.Y = y
}

func (c *Object) SetRotation(r CafeObjectRotation) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.rotation = r
}

func (c *Object) SetDishID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishID = id
}

func (c *Object) SetDishStatus(status int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishStatus = status
}

func (c *Object) SetDishAmount(amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishAmount = amount
}

func (c *Object) AddDishAmount(amount int) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.dishAmount+amount < 0 {
		return false
	}
	c.dishAmount += amount

	if c.dishAmount == 0 {
		c.dishID = -1
	}
	return true
}

func (c *Object) SetFancyIng(fancying bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.fancyIng = fancying
}

func (c *Object) SetStartedAt(t *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.startedAt = t
}

func (c *Object) SetFinishesAt(t *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.finishesAt = t
}

func (c *Object) SetOccupied(b bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.occupied = b
}
