package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"log"
)

func HandleClient(c *client.Client, gm *managers.GameManager) {

	for c.Alive() {

		// Read next request
		req, err := c.NextRequest()
		if err != nil {
			log.Printf("Failed to read request: %s\n", err.Error())
			break
		}

		// Handle requests
		err = HandleRequest(req, c, gm)
		if err != nil {
			log.Printf("Error while handling request: %s\n", err.Error())
			continue
		}

	}

	c.Disconnect()
}

func HandleRequest(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

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
		RoomList(req, c, gm)
	case requests.C2S_VERSION_CHECK:
		err = VersionCheck(req, c, gm)
	case requests.C2S_JOIN_CAFE:
		err = JoinCafe(req, c, gm)
	case requests.C2S_LOGIN:
		err = Login(req, c, gm)
	case requests.C2S_CAFE_WALK:
		err = CafeWalk(req, c, gm)
	case requests.C2S_SHOP_AVAILIBILITY:
		err = ShopAvailibility(req, c, gm)
	case requests.C2S_SHOP_DELETE_ITEM:
		err = SellIngredient(req, c, gm)
	case requests.C2S_SHOP_BUY_ITEM:
		err = BuyIngredient(req, c, gm)
	case requests.C2S_CAFE_CHAT:
		err = SendChatMessage(req, c, gm)
	case requests.C2S_MARKETPLACE_JOIN:
		err = JoinMarketplace(req, c, gm)
	case requests.C2S_CAFE_COOK:
		err = StartCooking(req, c, gm)
	case requests.C2S_CAFE_STOVE_DELIVER_INFO:
		err = StoveDeliverInfo(req, c, gm)
	case requests.C2S_CAFE_STOVE_DELIVER:
		err = StoveDeliver(req, c, gm)
	case requests.C2S_CAFE_CLEAN:
		err = Clean(req, c, gm)
	case requests.C2S_EDITOR_MODE:
		err = EditorMode(req, c, gm)
	case requests.C2S_EDITOR_BUY_OBJECT:
		// TODO: Need to check level.
		err = BuyObject(req, c, gm)
	case requests.C2S_EDITOR_STORE_OBJECT:
		err = StoreObject(req, c, gm)
	case requests.C2S_EDITOR_ROTATE_OBJECT:
		err = RotateObject(req, c, gm)
	case requests.C2S_EDITOR_MOVE_OBJECT:
		err = MoveObject(req, c, gm)
	case requests.C2S_EDITOR_SELL_OBJECT:
		err = SellObject(req, c, gm)
	case requests.C2S_CAFE_INSTANTCOOK:
		err = InstantCook(req, c, gm)
	default:
		log.Printf("[NOT IMPLEMENTED]: %v\n", req.Args[0])
	}

	return err
}
