package coops

import (
	"cafego/internal/models/balancing"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Coop struct {
	ID            int             `gorm:"primaryKey;autoIncrement;column:id"`
	ActiveCoop    int             `gorm:"column:active_coop;not null"`
	Host          int             `gorm:"column:host;not null"`
	Members       simple.IntSlice `gorm:"column:members;type:text;not null"`
	Dishes        simple.IntMap   `gorm:"column:dishes;type:text;not null"`
	Start         time.Time       `gorm:"column:start;default:CURRENT_TIMESTAMP"`
	End           time.Time       `gorm:"column:end;not null"`
	FinishLevel   int             `gorm:"column:finish_level;not null"`
	ExtendCount   int             `gorm:"column:extend_count;default: 0"`
	runtime       int             `gorm:"-"`
	playersString string          `gorm:"-"`
}

func (coop Coop) TableName() string {
	return "coop"
}

func (coop Coop) AsResponse() string {
	dishes := make([]string, 0, len(coop.Dishes))
	for dishID, amount := range coop.Dishes {
		dishes = append(dishes, fmt.Sprintf("%d+%d", dishID, amount))
	}
	dishesStr := strings.Join(dishes, "#")

	args := fmt.Sprintf("%d+%d+%d+%d", coop.ID, coop.ActiveCoop, coop.GetRuntime(), coop.ExtendCount)
	//     _loc3_ = ID+ACT+RUN+EXT
	return args + "%" + strconv.Itoa(coop.FinishLevel) + "%" + dishesStr + "%" + coop.playersString
}

func (coop Coop) AsActiveListResponse() string {
	args := fmt.Sprintf("%d+%d+%d+%d", coop.ID, coop.ActiveCoop, coop.GetRuntime(), coop.ExtendCount)
	return args
}

func (coop Coop) GetCoop() Coop {
	return coop
}

func (coop *Coop) GetRuntime() int {
	return int(time.Since(coop.Start).Seconds())
}

func (coop *Coop) GetIsActive() bool {
	if coop.FinishLevel != -1 {
		return false
	}
	return int(time.Since(coop.Start).Seconds()) <= coop.GetRuntime()
}

func (coop *Coop) SetPlayersString(playersString string) {
	coop.playersString = playersString
}

func (coop *Coop) CalculateFinishLevel() int {
	coopInfo, err := utils.GetCoop(coop.ActiveCoop)
	if err != nil {
		return -1
	}

	timeLeft := max(0, coopInfo.Duration+coop.ExtendCount*balancing.BalancingConstants.CoopExpansionHours*3600)
	runtime := coop.GetRuntime()

	timeLeftSilver := max(0, int(float64(timeLeft)*balancing.BalancingConstants.CoopTimeToSilver)-runtime)
	timeLeftGold := max(0, int(float64(timeLeft)*balancing.BalancingConstants.CoopTimeToGold)-runtime)

	if timeLeftGold > 0 {
		return 0
	}

	if timeLeftSilver > 0 {
		return 1
	}

	if timeLeft >= 0 {
		return 2
	}

	return -1
}

func (coop *Coop) IsDone() bool {
	coopInfo, _ := utils.GetCoop(coop.ActiveCoop)
	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		dishID, _ := strconv.Atoi(dishRequirement[0])
		dishAmount, _ := strconv.Atoi(dishRequirement[1])

		if coop.Dishes[dishID] == dishAmount {
			continue
		} else {
			return false
		}
	}
	return true
}

func (coop *Coop) SetHost(playerID int) {
	coop.Host = playerID
}

func (coop *Coop) Join(playerID int) {
	coop.Members = append(coop.Members, playerID)
}

func (coop *Coop) AddDish(dishID int) {
	coopInfo, _ := utils.GetCoop(coop.ActiveCoop)

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		reqDishID, _ := strconv.Atoi(dishRequirement[0])
		reqDishAmount, _ := strconv.Atoi(dishRequirement[1])

		if reqDishID == dishID {
			if coop.Dishes[dishID] < reqDishAmount {
				coop.Dishes[dishID]++
			}
		}
	}
}

func (coop *Coop) Leave(playerID int) {
	if len(coop.Members) != 0 {
		newMembers := make(simple.IntSlice, 0, len(coop.Members))
		for _, memberID := range coop.Members {
			if memberID != playerID {
				newMembers = append(newMembers, memberID)
			}
		}

		coop.Members = newMembers

		if coop.Host == playerID {
			for _, memberID := range coop.Members {
				coop.Host = memberID
				break
			}
		}
	}
}
