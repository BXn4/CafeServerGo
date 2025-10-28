package managers

import (
	"cafego/internal/utils"
	"fmt"
	"math/rand"
)

func (gm *GameManager) SendEarnAchievement(achievementWOD int, level int, username string) {
	c, err := gm.GetClientByName(username)
	if err != nil {
		return
	}

	achievementID := utils.GetAchievementIDByWOD(achievementWOD)
	achievementXPReward := utils.GetAchievementXPReward(achievementWOD, level-1)     // to get previous reward
	achievementCashReward := utils.GetAchievementCashReward(achievementWOD, level-1) //  to get previous reward
	achievementGoldReward := utils.GetAchievementGoldReward(achievementWOD, level-1) //  to get previous reward

	args := fmt.Sprintf("%d+%d", achievementID, level)

	c.Player.AddXP(achievementXPReward)
	c.Player.AddCash(achievementCashReward)
	c.Player.AddGold(achievementGoldReward)

	// 2005 -> boughtingredients
	// 2013 -> coopgoldcount
	// 2014 -> 2019 total dish category cooked
	// 2020 -> total dishes cooked
	// 2021 -> different dishes
	// 2022 -> masteries
	// 2023 -> masteries gold
	// 2024 -> fancy count
	if achievementWOD == 2005 || achievementID >= 2013 && achievementID <= 2024 {
		fancies, err := utils.GetItems("fancy")
		if err != nil {
			return
		}
		choice := rand.Intn(len(fancies))
		fancy := fancies[choice]
		amount := 10

		args = fmt.Sprintf("%s+%d+%d", args, fancy.ID, amount)

		c.Location.Cafe().AddToFridge(fancy.ID, amount)

		c.DB.UpdateFridgeInventory(c.Location.Cafe().ID, c.Location.Cafe().FridgeInventory.String())
	}

	c.DB.UpdateXP(c.Player.ID, c.Player.GetXP())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

	c.SendExtensionResponse("cae", "-1", "0", args)
}
