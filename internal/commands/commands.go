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
	if req.Kind == requests.POLICY_FILE {
		c.SendSystemResponse(responses.POLICY_FILE)

	} else if req.Kind == requests.VERSION_CHECK {
		c.SendSystemResponse(responses.VERSION_CHECK)

	} else if req.Kind == requests.AUTO_JOIN {
		// TODO: Handle
		c.SendSystemResponse(responses.AUTO_JOIN)

	} else if req.Kind == requests.ROUND_TRIP {
		c.SendSystemResponse(responses.ROUND_TRIP)

	} else if req.Kind == requests.DISCONNECT {
		c.SendSystemResponse(responses.LOGOUT)

	} else if req.Kind == requests.LOGIN {
		RoomList(req, c, clientManager, cafeManager)

	} else if req.Kind == requests.C2S_VERSION_CHECK {
		err = VersionCheck(req, c, clientManager, cafeManager)

	} else if req.Kind == requests.C2S_JOIN_CAFE {
		err = JoinCafe(req, c, clientManager, cafeManager)

	} else if req.Kind == requests.C2S_LOGIN {
		err = Login(req, c, clientManager, cafeManager) // lgn - S2C_LOGIN
		if err != nil {
			return err
		}

		//.....
		//}else if(req.Kind == ROUND_TRIP_REQUEST) {

	}

	if err != nil {
		return err
	}

	return nil
}
