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
		c.Player.DailyLogin = time.Now().UTC()
		c.DB.UpdateDailyLogin(c.Player.ID, c.Player.DailyLogin)
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

		c.Player.Cash += loginBonus
		c.Location.Cafe().AddToFridge(fancy.ID, 1)
		// c.DB.UpdateFridge()

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
		dishID := counter.GetDishID()
		if dishID == -1 {
			break
		}

		// If cant sell more continue
		if !counter.AddDishAmount(-1) {
			continue
		}
		soldDishes++

		// Get dish info
		wod, err := utils.GetDish(dishID)
		if err != nil {
			return err
		}

		// Add to passive income
		passiveIncome += wod.IncomePerServing
	}

	if isDaily || soldDishes > 0 {
		if soldDishes > 0 {
			c.Player.UpdateAchivementServingsCount(soldDishes)
		}
		soldDishesStr := strconv.Itoa(soldDishes)

		c.Player.AddCash(passiveIncome)

		// Add passive income to args
		args = append(args, fmt.Sprintf("1901+%d", passiveIncome))

		// Save modified cafe
		c.DB.UpdateObjects(mycafe.ID, mycafe.Objects.StringForDB())
		c.DB.UpdateCash(c.Player.ID, c.Player.Cash)
		c.DB.UpdateAchievement(c.Player.ID, c.Player.GetAchivements().String())

		// Send response
		c.SendExtensionResponse("lbu", "-1", "0", strings.Join(args, "#"), soldDishesStr)

		AssetsSynchronize(c) // Updates force the player cash, gold in the game visually.
		// Its used in the payments, but we can use it to update the cash, gold values force.
	}

	return nil
}
