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

				toClient.DB.UpdateCoopID(toClient.Player.ID, toClient.Player.GetActiveCoopID())
			}
		}
	}

	return nil
}
