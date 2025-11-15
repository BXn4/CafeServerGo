package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/coops"
	"strconv"
)

func CoopFinish(coop *coops.Coop, gm *managers.GameManager) error {
	for _, memberID := range coop.Members {
		item, err := gm.GetClient(memberID)
		var toClient *client.Client
		if err == nil {
			toClient = item.(*client.Client)
			if toClient.Player.GetActiveCoopID() == coop.ID {
				toClient.SendExtensionResponse("cof", "-1", strconv.Itoa(coop.FinishLevel))
				toClient.Player.SetActiveCoopID(0)

				gold, cash, xp := coop.GetRewards(coop.FinishLevel, toClient.Player.GetLevel())

				toClient.Player.AddGold(gold)
				toClient.Player.AddCash(cash)
				toClient.Player.AddXP(xp)

				toClient.DB.UpdateGold(toClient.Player.ID, toClient.Player.GetGold())
				toClient.DB.UpdateCash(toClient.Player.ID, toClient.Player.GetActiveCoopID())
				toClient.DB.UpdateXP(toClient.Player.ID, toClient.Player.GetXP())

				toClient.DB.UpdateCoopID(toClient.Player.ID, toClient.Player.GetActiveCoopID())
			}
		}
	}

	return nil
}
