package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
)

// min level 4

/* func init() {
	RegisterCommand(requests.C2S_MARKETPLACE_JOB,
		CommandConfig{
			Name:       "MarketplaceJob",
			Identifier: responses.S2C_MARKETPLACE_JOB,
			MinArgs:    6,
			MaxArgs:    6,
		},
		AllowFriendRequestsValidator,
		AllowFriendRequests,
	)
}

const (
	OFFER_JOB                  = 0
	ACCEPT_JOB                 = 1
	DECLINE_JOB                = 2
	PLAYER_LEFT                = 11
	PLAYER_WORKED_ENOUGH_TODAY = 61
) */

// the func is "done" but the job is not yet.
func MarketplaceJob(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	/*
		status, err := strconv.Atoi(req.Args[2])
		fromPlayerID, err := strconv.Atoi(req.Args[3])
		toPlayerID, err := strconv.Atoi(req.Args[4])

		if err != nil {
			return fmt.Errorf("Cant convert string to int!")
		}

		item, err := gm.GetClient(fromPlayerID)
		var fromClient *client.Client
		// var _ *player.Player
		if err == nil {
			fromClient = item.(*client.Client)
			// fromPlayer = fromClient.Player
		}

		item, err = gm.GetClient(toPlayerID)
		var toClient *client.Client
		var toPlayer *player.Player
		if err == nil {
			toClient = item.(*client.Client)
			toPlayer = toClient.Player
		}

		if toPlayer == nil {
			status = PLAYER_LEFT
		}

		if toPlayer != nil && toClient.Location != fromClient.Location {
			status = PLAYER_LEFT
		}

		if toPlayer.OpenJobs < 1 {
			status = PLAYER_WORKED_ENOUGH_TODAY
		}

		switch status {
		case OFFER_JOB:
			// To dont allow spam offers, because more players can offer jobb offers, and the same player can send multiple times. dont allow it
			if !slices.Contains(toPlayer.Job.Offers, fromPlayerID) {
				toPlayer.AddJobOffer(fromPlayerID)
				toClient.SendExtensionResponse("mjo", "-1", "0", strconv.Itoa(status), strconv.Itoa(fromPlayerID), strconv.Itoa(toPlayerID))
			}
			return nil
		case ACCEPT_JOB:
			if slices.Contains(toClient.Player.Job.Offers, fromPlayerID) {
				toClient.Player.ClearOffers()
				fromClient.SendExtensionResponse("mjo", "-1", "0", strconv.Itoa(status), strconv.Itoa(fromPlayerID), strconv.Itoa(toPlayerID))

				toClient.SendExtensionResponse("mjo", "-1", "0", strconv.Itoa(status), strconv.Itoa(toPlayerID), strconv.Itoa(toPlayerID))

				toPlayer.StartJob(fromClient.Location.Cafe())

				toPlayer.OpenJobs--

				toClient.DB.UpdateOpenJobs(toPlayer.ID, toPlayer.OpenJobs)

				JoinCafe(req, toClient, gm)
			}
			return nil
		case DECLINE_JOB:
			fromClient.SendExtensionResponse("mjo", "-1", "0", strconv.Itoa(status), strconv.Itoa(fromPlayerID), strconv.Itoa(toPlayerID))

			toClient.Player.RemoveJobOffer(fromPlayerID)
			return nil
		case PLAYER_LEFT:
			// if(_loc5_ == CafeModel.userData.userID)
			// toClient.SendExtensionResponse("mjo", "-1", "0", strconv.Itoa(status), strconv.Itoa(toPlayerID), strconv.Itoa(fromPlayerID))

			toClient.Player.RemoveJobOffer(fromPlayerID)
			return nil
		case PLAYER_WORKED_ENOUGH_TODAY:
			fromClient.SendExtensionResponse("mjo", "-1", strconv.Itoa(status), strconv.Itoa(fromPlayerID), strconv.Itoa(fromPlayerID))
			return nil
			} */
	return nil
}
