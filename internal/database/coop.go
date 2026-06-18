package database

import (
	"cafego/internal/models/coops"
	"cafego/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CreateCoop creates a new cooperative game session
func (db *CafeDB) CreateCoop(coopID, playerID int, end time.Time) (int, error) {
	coop := coops.Coop{
		ActiveCoop:  coopID,
		Host:        playerID,
		Members:     []int{playerID}, // Initially only contains the host
		Dishes:      make(map[int]int),
		Start:       time.Now().UTC(),
		End:         end,
		FinishLevel: -1,
	}

	coopInfo, err := utils.GetCoop(coopID)
	if err != nil {
		return 0, fmt.Errorf("Cannot get dishinfo")
	}

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		dishID, _ := strconv.Atoi(dishRequirement[0])

		coop.Dishes[dishID] = 0
	}

	result := db.conn.Create(&coop)
	if result.Error != nil {
		return 0, result.Error
	}

	return coop.ID, nil
}

// GetCoop retrieves a cooperative game session by ID
func (db *CafeDB) GetCoop(id int) (coops.Coop, error) {
	var c coops.Coop
	result := db.conn.Where("id = ?", id).First(&c)
	if result.Error != nil {
		return coops.Coop{}, result.Error
	}

	var playersList []string

	for _, playerID := range c.Members {
		player, err := db.GetPlayer(playerID)

		if err != nil {
			playersList = append(playersList, "")
		} else {
			avatar := player.GetAvatar()
			playersList = append(playersList, strconv.Itoa(player.GetID())+"+"+strconv.Itoa(player.GetXP())+"+"+avatar.String())
		}
	}

	c.SetPlayersString(strings.Join(playersList, "%"))

	return c, nil
}

func (db *CafeDB) GetCoopByHost(playerID int) (coops.Coop, error) {
	var c coops.Coop
	result := db.conn.Where("host = ?", playerID).First(&c)
	if result.Error != nil {
		return coops.Coop{}, result.Error
	}

	/* playersString := make([]string, len(c.Members))

	for i, playerID := range c.Members {
		player, err := db.GetPlayer(playerID)

		if err != nil {
			playersString[i] = ""
		} else {
			playersString[i] = strconv.Itoa(player.ID) + "+" + strconv.Itoa(player.XP) + "+" + player.Avatar.String()
		}
	}

	c.SetPlayersString(playersString) */
	return c, nil
}

// DeleteCoop removes a cooperative game session by ID
func (db *CafeDB) DeleteCoop(id int) error {
	result := db.conn.Delete(&coops.Coop{}, id)
	return result.Error
}

func (db *CafeDB) SaveCoop(coop *coops.Coop) error {
	result := db.conn.Save(&coop)
	return result.Error
}

/*
// AddMemberToCoop adds a player to an existing coop session
func (db *CafeDB) AddMemberToCoop(coopID int, playerID int) error {
	coop, err := db.GetCoop(coopID)
	if err != nil {
		return err
	}

	coop.Members = append(coop.Members, playerID)
	result := db.conn.Save(&coop)
	return result.Error
}

func (db *CafeDB) RemoveMemberFromCoop(coopID, playerID int) error {
	coop, err := db.GetCoop(coopID)
	if err != nil {
		return err
	}

	membersRebuild := make([]int, len(coop.Members))
	for memberID := range coop.Members {
		if playerID != memberID {
			membersRebuild = append(membersRebuild, memberID)
		}
	}

	coop.Members = membersRebuild

	result := db.conn.Save(&coop)
	return result.Error
}

// UpdateCoopDishes updates the dishes count in a coop session
func (db *CafeDB) UpdateCoopDishes(coopID int, dishID int, count int) error {
	coop, err := db.GetCoop(coopID)
	if err != nil {
		return err
	}

	coopInfo, _ := utils.GetCoop(coop.ActiveCoop)

	coopDishes := strings.Split(coopInfo.Dishes, "#")

	for _, dishRequirements := range coopDishes {
		dishRequirement := strings.Split(dishRequirements, "+")
		reqDishID, _ := strconv.Atoi(dishRequirement[0])
		reqDishAmount, _ := strconv.Atoi(dishRequirement[1])

		if reqDishID == dishID {
			if reqDishAmount != coop.Dishes[dishID] {
				coop.Dishes[dishID]++
			}
		}
	}

	result := db.conn.Save(&coop)
	return result.Error
}
*/
