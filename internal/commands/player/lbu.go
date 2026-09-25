/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package player

import (
	"cafego/internal/agents"
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func LoginRewards(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	var args []string

	// Get my cafe
	cafe, err := c.DB.GetCafeByPlayerID(c.Player.GetID())

	if err != nil {
		return err
	}

	if c.Player.GetIsDailyLogin() {

		c.Player.SetDailyLogin(time.Now().UTC())

		maxInstants := utils.GetLevelInstantCookingsLimit(c.Player.GetLevel())

		c.Player.SetMaxInstants(maxInstants)

		c.Player.SetInstantCookings(0) // 0 because its counts used, and not how much still have

		c.Player.SetPlayedWheel(false)

		if c.Player.GetCoopID() == 0 {
			c.Player.SetIsStartedCoop(false)

			c.DB.UpdateStartedCoop(c.Player.GetID(), c.Player.GetIsStartedCoop())

		}

		c.Player.SetOpenJobs(5)
		fancies, err := utils.GetItems("fancy")

		if err != nil {
			return err
		}

		choice := rand.Intn(len(fancies))
		fancy := fancies[choice]
		amount := 1

		// Calculate login bonus
		loginBonus := 500 * (int(c.Player.GetLevel()/10) + 1)
		dailyLoginBonusStr := fmt.Sprintf("1906+%d#%d+%d", loginBonus, fancy.ID, amount)
		args = append(args, dailyLoginBonusStr)

		c.Player.AddCash(loginBonus)
		cafe.AddToFridge(fancy.ID, amount)

		c.DB.UpdateDailyLogin(c.Player.GetID(), c.Player.GetDailyLogin())
		c.DB.UpdatePlayedWheel(c.Player.GetID(), c.Player.GetPlayedWheel())
		c.DB.UpdateInstantCookings(c.Player.GetID(), maxInstants)
		c.DB.UpdateOpenJobs(c.Player.GetID(), c.Player.GetOpenJobs())
		c.DB.UpdateFridgeInventory(cafe.GetID(), cafe.GetFridgeInventory().String())
	}

	canSellDishes := false
	soldDishes := 0

	for _, obj := range cafe.GetObjects() {
		if obj.IsCounter() {
			if obj.GetDishAmount() > 0 {
				canSellDishes = true
				break
			}
		}
	}

	if canSellDishes {
		rating := cafe.GetRating()
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

		secondsPassedSinceLastLogin := time.Now().UTC().Sub(c.Player.GetLastLogin()).Seconds()
		maxShouldSellDishCount := int(secondsPassedSinceLastLogin / float64(customerSpawnTime))

		passiveIncome := 0

		for i := 0; i < maxShouldSellDishCount; i++ {
			// Get counter
			counter, _ := agents.GetRandomCounter(cafe)

			if counter == nil {
				break
			}

			dishID := counter.GetDishID()
			dishAmount := counter.GetDishAmount()

			if dishID == -1 || dishAmount <= 0 {
				break
			}

			if !counter.AddDishAmount(-1) {
				continue
			}

			soldDishes++

			wod, err := utils.GetDish(dishID)
			if err != nil {
				return err
			}

			passiveIncome += wod.IncomePerServing
		}

		if soldDishes > 0 {
			dishesSoldOffline := fmt.Sprintf("1901+%d", passiveIncome)
			args = append(args, dishesSoldOffline)

			c.Player.AddCash(passiveIncome)

			c.Player.UpdateAchivementServingsCount(soldDishes)
			c.DB.UpdateObjects(cafe.GetID(), cafe.GetObjects().StringForDB())
			c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
		}
	}

	if len(args) > 0 {
		c.SendExtensionResponse(responses.S2C_LOGIN_BONUS, "-1", "0", strings.Join(args, "#"), strconv.Itoa(soldDishes))

		c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
		c.DB.UpdateAchievement(c.Player.GetID(), c.Player.GetAchivements().String())
		c.DB.UpdateCash(c.Player.GetID(), c.Player.GetCash())
	}

	return nil
}
