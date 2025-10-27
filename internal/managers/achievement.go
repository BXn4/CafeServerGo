package managers

import (
	"cafego/internal/utils"
	"fmt"
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

	c.DB.UpdateXP(c.Player.ID, c.Player.GetXP())
	c.DB.UpdateCash(c.Player.ID, c.Player.GetCash())
	c.DB.UpdateGold(c.Player.ID, c.Player.GetGold())

	c.SendExtensionResponse("cae", "-1", "0", args)
}
