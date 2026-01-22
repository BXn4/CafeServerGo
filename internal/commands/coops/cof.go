package coops

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/coops"
	"cafego/internal/types/responses"
	"strconv"
)

func CoopFinish(coop *coops.Coop, gm *managers.GameManager) error {
	for _, memberID := range coop.Members {
		item, err := gm.GetClient(memberID)
		var toClient *client.Client
		if err == nil {
			toClient = item.(*client.Client)
			if toClient.Player.GetCoopID() == coop.ID {
				toClient.SendExtensionResponse(responses.S2C_COOP_FINISH, "-1", strconv.Itoa(coop.FinishLevel))
				toClient.Player.SetCoopID(0)

				gold, cash, xp := coop.GetRewards(coop.FinishLevel, toClient.Player.GetLevel())

				toClient.Player.AddGold(gold)
				toClient.Player.AddCash(cash)
				toClient.Player.AddXP(xp)

				toClient.DB.UpdateGold(toClient.Player.GetID(), toClient.Player.GetGold())
				toClient.DB.UpdateCash(toClient.Player.GetID(), toClient.Player.GetCoopID())
				toClient.DB.UpdateXP(toClient.Player.GetID(), toClient.Player.GetXP())

				toClient.DB.UpdateCoopID(toClient.Player.GetID(), toClient.Player.GetCoopID())
			}
		}
	}

	return nil
}
