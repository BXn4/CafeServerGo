package commands

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func LoginRewards(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	// Check if time passed by daily login
	timePassedSinceLastLogin := time.Now().UTC().Sub(c.Player.LastLogin)
	timePassed := time.Now().UTC().Sub(c.Player.DailyLogin)
	isDaily := timePassed >= 24*time.Hour

	if isDaily {
		c.Player.DailyLogin = time.Now()
	}

	var args []string
	loginBonusStr := ""

	if isDaily {
		// Get random fancy
		fancies, err := utils.GetItems("fancy")
		if err != nil {
			return err
		}
		choice := rand.Intn(len(fancies))
		fancy := fancies[choice]

		// Calculate login bonus
		loginBonus := 500 * (int(c.Player.GetLevel()/10) + 1)
		loginBonusStr = fmt.Sprintf("1906+%v#%v+1", loginBonus, fancy.ID)
		args = append(args, loginBonusStr)

	}

	// Get my cafe
	mycafe, err := c.DB.GetCafeByPlayerID(c.Player.ID)

	if err != nil {
		return err
	}

	// Calculate customer spawn time (we use max time so people dont try to cheat the system)
	rating := mycafe.GetRating()
	var customerSpawnTime int

	if rating < 150 {
		customerSpawnTime = 30
	} else if rating <= 150 && rating < 350 {
		customerSpawnTime = 20
	} else if rating <= 350 && rating < 500 {
		customerSpawnTime = 15
	} else {
		customerSpawnTime = 5
	}

	// Calculate passive income
	passedSeconds := timePassedSinceLastLogin.Seconds()
	passiveIncome := 0
	soldDishes := 0
	maxDishSellCount := int(passedSeconds / float64(customerSpawnTime))

	for i := 0; i < maxDishSellCount; i++ {
		// Get counter
		counter, _ := agents.GetRandomCounter(mycafe)

		// If cant find counter break
		if counter == nil {
			break
		}

		// If no food return
		savedID := counter.GetDishID()
		if savedID <= 0 {
			break
		}

		// If cant sell more continue
		if !counter.AddDishAmount(-1) {
			continue
		}
		soldDishes++

		// Get dish info
		wod, err := utils.GetDish(savedID)
		if err != nil {
			return err
		}

		// Add to passive income
		passiveIncome += wod.Cash
	}

	if soldDishes > 0 {

		soldDishesStr := strconv.Itoa(soldDishes)

		c.Player.Cash += passiveIncome

		// Add passive income to args
		args = append(args, fmt.Sprintf("1901+%v", passiveIncome))

		// Save modified cafe
		c.DB.UpdateObjects(mycafe.ID, mycafe.Objects.StringForDB())
		c.DB.UpdateCash(c.Player.ID, c.Player.Cash)

		// Send response
		c.SendExtensionResponse("lbu", "-1", "0", strings.Join(args, "#"), soldDishesStr)
	}

	return nil
}
