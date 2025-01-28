package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/objects"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

const (
	REQUEST  = 0
	ACCEPT   = 1
	DENY     = 2
	UNFRIEND = 3
)

func SendFriendRequest(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	action, err := strconv.Atoi(req.Args[2]) // REQUEST, ACCEPT, DENY, UNFRIEND
	if err != nil {
		return err
	}

	fromPlayerID, err := strconv.Atoi(req.Args[3])
	if err != nil {
		return err
	}

	toPlayerID, err := strconv.Atoi(req.Args[4])
	if err != nil {
		return err
	}

	// Get from player
	item, err := gm.GetClient(fromPlayerID)
	var fromClient *client.Client
	var fromPlayer *objects.Player
	if err == nil {
		fromClient = item.(*client.Client)
		fromPlayer = fromClient.Player
	}

	// Get to player
	item, err = gm.GetClient(toPlayerID)
	var toClient *client.Client
	var toPlayer *objects.Player
	if err == nil {
		toClient = item.(*client.Client)
		toPlayer = toClient.Player
	}

	// Handle actions
	switch action {
	case REQUEST:
		if toPlayer == nil {
			break
		}

		if toPlayer.AllowFriendRequests {
			toClient.SendExtensionResponse("big", "-1", "0", "0", req.Args[3], req.Args[4])
		}
	case ACCEPT:
		if toPlayer == nil {
			break
		}

		// Add friends
		toPlayer.AddFriend(fromPlayer.ID)
		fromPlayer.AddFriend(toPlayer.ID)

		// Send messages
		err = SendFriendsAvatar(nil, toClient, gm)
		if err != nil {
			return err
		}

		err = SendFriendsAvatar(nil, fromClient, gm)
		if err != nil {
			return err
		}
	case DENY:
		if toPlayer == nil {
			break
		}
		fromClient.SendExtensionResponse("big", "-1", "0", "2",
			req.Args[3],
			req.Args[4],
			fmt.Sprintf("%v+%v+%v", toPlayer.ID, toPlayer.XP, toPlayer.Avatar.String(toPlayer.Username)),
		)
	case UNFRIEND:
		// TODO: Test

		fromPlayer.DeleteFriend(toPlayer.ID)

		if toPlayer != nil {
			// if player online
			toPlayer.DeleteFriend(fromPlayer.ID)
			fromClient.SendExtensionResponse("big", "-1", "0", "3",
				req.Args[3],
				req.Args[4],
				fmt.Sprintf("%v+%v+%v", toPlayer.ID, toPlayer.XP, toPlayer.Avatar.String(toPlayer.Username)),
			)
		} else {
			// if player offline
			c.DB.DeleteFriend(toPlayerID, fromPlayer.ID)
		}

	}

	return nil
}
