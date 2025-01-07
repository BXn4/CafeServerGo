package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	_ "errors"
	"fmt"
)

func HandleClient(c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) {

	for c.Alive() {

		// Read next request
		req, err := c.NextRequest()
		if err != nil {
			fmt.Printf("Failed to read request: %s\n", err.Error())
			break
		}

		// Handle requests
		err = HandleRequest(req, c, clientManager, cafeManager)
		if err != nil {
			fmt.Printf("Failed to handle request: %s\n", err.Error())
			break
		}

	}

	c.Disconnect()
}

func HandleRequest(req *requests.Request, c *client.Client, clientManager *managers.ClientManager, cafeManager *managers.CafeManager) error {

	var err error

	switch req.Kind {
	/* SYSTEM */
	case requests.POLICY_FILE:
		c.SendSystemResponse(responses.POLICY_FILE)
	case requests.VERSION_CHECK:
		c.SendSystemResponse(responses.VERSION_CHECK)
	case requests.AUTO_JOIN:
		c.SendSystemResponse(responses.AUTO_JOIN)
	case requests.ROUND_TRIP:
		c.SendSystemResponse(responses.ROUND_TRIP)
	case requests.DISCONNECT:
		c.SendSystemResponse(responses.LOGOUT)

	/* COMMANDS */
	case requests.LOGIN:
		RoomList(req, c, clientManager, cafeManager)
	case requests.C2S_VERSION_CHECK:
		err = VersionCheck(req, c, clientManager, cafeManager)
	case requests.C2S_JOIN_CAFE:
		err = JoinCafe(req, c, clientManager, cafeManager)
	case requests.C2S_LOGIN:
		err = Login(req, c, clientManager, cafeManager)
		if err != nil {
			return err
		}
	case requests.C2S_CAFE_WALK:
		err = CafeWalk(req, c, clientManager, cafeManager)
	case requests.C2S_SHOP_AVAILIBILITY:
		err = ShopAvailibility(req, c, clientManager, cafeManager)
	}

	if err != nil {
		return err
	}

	return nil
}
