package objects

import (
	"cafego/internal/types/cafetypes"
	"cafego/internal/types/daos"
	"cafego/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type CafeObject struct {
	kind     cafetypes.CafeObjectKind
	pos      [2]int
	rotation cafetypes.CafeObjectRotation

	dishID     int
	dishStatus int
	dishAmount int

	fancyIng   bool
	startedAt  *time.Time
	finishesAt *time.Time

	mutex sync.Mutex
}

func NewCafeObject(posX int, posY int, objID int, objRotation int) (*CafeObject, error) {
	cafeObj := CafeObject{
		pos:      [2]int{posX, posY},
		kind:     cafetypes.CafeObjectKind(objID),
		rotation: cafetypes.CafeObjectRotation(objRotation),
	}
	return &cafeObj, nil
}

func NewCafeObjectFromJSON(s string) (*CafeObject, error) {

	var cafeObj CafeObject
	if err := json.Unmarshal([]byte(s), &cafeObj); err != nil {
		return nil, err
	}

	return &cafeObj, nil
}

func NewCafeObjectFromString(s string) (*CafeObject, error) {

	data := strings.Split(s, "+")
	items, err := utils.MultiAtoi(data...)
	if err != nil {
		return nil, fmt.Errorf("Error parsing cafeobject from string: %v", err)
	}

	posX, posY, kind, rotation := items[0], items[1], items[2], items[3]

	cafeObj := CafeObject{
		pos:      [2]int{posX, posY},
		kind:     cafetypes.CafeObjectKind(kind),
		rotation: cafetypes.CafeObjectRotation(rotation),
	}

	return &cafeObj, nil
}

func (c *CafeObject) IsWall() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 101 <= c.kind && c.kind <= 135
}

func (c *CafeObject) IsDoor() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 201 <= c.kind && c.kind <= 207
}

func (c *CafeObject) IsStove() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 252 <= c.kind && c.kind <= 259
}

func (c *CafeObject) isStove() bool {
	return 252 <= c.kind && c.kind <= 259
}

func (c *CafeObject) IsCounter() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 301 <= c.kind && c.kind <= 312
}

func (c *CafeObject) isCounter() bool {
	return 301 <= c.kind && c.kind <= 312
}

func (c *CafeObject) IsChair() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 601 <= c.kind && c.kind <= 624
}

func (c *CafeObject) isChair() bool {
	return 601 <= c.kind && c.kind <= 624
}

func (c *CafeObject) IsTable() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 401 <= c.kind && c.kind <= 423
}

func (c *CafeObject) IsVending() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.kind == 1701
}

func (c *CafeObject) isVending() bool {
	return c.kind == 1701
}

func (c *CafeObject) IsFridge() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return 351 <= c.kind && c.kind <= 358
}

func (c *CafeObject) String() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	args := []string{
		strconv.Itoa(c.pos[0]),
		strconv.Itoa(c.pos[1]),
		strconv.Itoa(int(c.kind)),
		strconv.Itoa(int(c.rotation)),
	}

	// if STOVE
	if c.isStove() {

		args = append(args, strconv.Itoa(c.dishID))
		if !(c.dishID == -1 || c.dishID == -2) {

			fancyIngStr := utils.If(c.fancyIng, "1", "0")
			args = append(args, fancyIngStr)

			if c.startedAt != nil {
				currentTime := time.Now().UTC()
				passedTime := currentTime.Sub(*c.startedAt).Seconds()
				remainingTime := c.finishesAt.Sub(currentTime).Seconds()

				args = append(args,
					strconv.Itoa(int(passedTime)),
					strconv.Itoa(int(remainingTime)),
				)
			} else {
				args = append(args, "-1", "-1")
			}
		}
		// if COUNTER
	} else if c.isCounter() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishAmount))
		// if CHAIR
	} else if c.isChair() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishAmount))
		// if VENDING
	} else if c.isVending() {
		args = append(args, strconv.Itoa(c.dishID), strconv.Itoa(c.dishAmount))
	}

	return strings.Join(args, "+")
}

func (c *CafeObject) JSON() string {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	obj := daos.CafeObjectDAO{
		Kind:     c.kind,
		Pos:      c.pos,
		Rotation: c.rotation,

		DishID:     c.dishID,
		DishStatus: c.dishStatus,
		DishAmount: c.dishAmount,

		FancyIng:   c.fancyIng,
		StartedAt:  c.startedAt,
		FinishesAt: c.finishesAt,
	}

	b, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("%v", err)
		return ""
	}
	return string(b)
}

func (c *CafeObject) GetNormalizedRotation() [2]int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.rotation == cafetypes.Up {
		return [2]int{1, 0}
	} else if c.rotation == cafetypes.Down {
		return [2]int{-1, 0}
	} else if c.rotation == cafetypes.Left {
		return [2]int{0, -1}
	} else /* Right */ {
		return [2]int{0, 1}
	}
}

// --- GETTERS -------------------------------------

func (c *CafeObject) GetKind() cafetypes.CafeObjectKind {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.kind
}

func (c *CafeObject) GetPos() [2]int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.pos
}

func (c *CafeObject) GetRotation() cafetypes.CafeObjectRotation {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.rotation
}

func (c *CafeObject) GetDishID() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishID
}

func (c *CafeObject) GetDishStatus() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishStatus
}

func (c *CafeObject) GetDishAmount() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.dishAmount
}

func (c *CafeObject) GetFancyIng() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.fancyIng
}

func (c *CafeObject) GetStartedAt() *time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.startedAt
}

func (c *CafeObject) GetFinishesAt() *time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.finishesAt
}

// --- SETTERS -------------------------------------

func (c *CafeObject) SetKind(k cafetypes.CafeObjectKind) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.kind = k
}

func (c *CafeObject) SetPos(pos [2]int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pos = pos
}

func (c *CafeObject) SetPosXY(x, y int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pos[0] = x
	c.pos[1] = y
}

func (c *CafeObject) SetRotation(r cafetypes.CafeObjectRotation) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.rotation = r
}

func (c *CafeObject) SetDishID(id int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishID = id
}

func (c *CafeObject) SetDishStatus(status int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishStatus = status
}

func (c *CafeObject) SetDishAmount(amount int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.dishAmount = amount
}

func (c *CafeObject) AddDishAmount(amount int) bool {
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

func (c *CafeObject) SetFancyIng(fancying bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.fancyIng = fancying
}

func (c *CafeObject) SetStartedAt(t *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.startedAt = t
}

func (c *CafeObject) SetFinishesAt(t *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.finishesAt = t
}
