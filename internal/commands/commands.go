package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"

	"github.com/charmbracelet/log"
)

func HandleClient(c *client.Client, gm *managers.GameManager) {
	defer c.Disconnect()
	for req := range c.RequestQueue {
		if req == nil {
			return
		}

		if req.NeedsLogin() && c.Player == nil &&
			req.Kind != requests.C2S_LOGIN && req.Kind != requests.C2S_SPECIAL_EVENT {
			// While the player is not logged in, disconnects the client if the request is not for login.
			// The client sends SEE command after login. Maybe we can patch this in the game client to only send it, when the login was successful.
			return
		}

		// Handle requests
		err := HandleRequest(req, c, gm)
		if err != nil {
			log.Warnf("%v request: %s", req.Args[0], err.Error())
			continue
		}
	}

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
	case requests.C2S_PING:
		SendPing(req, c, gm)
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
	case requests.C2S_MINI_MUFFIN:
		// TODO: Need to check level.
		err = PlayMuffinGame(req, c, gm)
	case requests.C2S_NPC_HIRE:
		err = HireWaiter(req, c, gm)
	case requests.C2S_NPC_FIRE:
		err = FireWaiter(req, c, gm)
	case requests.C2S_CAFE_ACHIEVEMENT_LIST:
		err = SendAchivements(req, c, gm)
	case requests.C2S_NPC_CUSTOMIZE:
		err = WaiterCustomize(req, c, gm)
	case requests.C2S_BUDDY_INGAME:
		err = SendFriendRequest(req, c, gm)
	case requests.C2S_EDITOR_BUY_FLOOR:
		err = BuyFloor(req, c, gm)
	case requests.C2S_WHEELOFFORTUNE:
		err = WheelOfFortune(req, c, gm)
	case requests.C2S_GIFT_PLAYERGIFTS:
		err = SendPlayerGifts(req, c, gm)
	case requests.C2S_GIFT_REMOVE:
		err = RemoveGift(req, c, gm)
	case requests.C2S_GIFT_USE:
		err = UseGift(req, c, gm)
	case requests.C2S_CREATE_AVATAR:
		err = CreateAvatar(req, c, gm)
	case requests.C2S_REGISTER:
		err = Register(req, c, gm)
	case requests.C2S_GIFT_SENDABLEGIFTS:
		err = DailyGifts(req, c, gm)
	case requests.C2S_GIFT_ALLREADYSEND_PLAYERS:
		err = GiftAllReadySendPlayers(req, c, gm)
	case requests.C2S_ALLOW_BUDDY_REQUESTS:
		err = AllowFriendRequests(req, c, gm)
	case requests.C2S_KICK_USER:
		err = KickPlayer(req, c, gm)
	case requests.C2S_ALLOW_MAIL_REQUESTS:
		err = AllowEmails(req, c, gm)
	case requests.C2S_CHANGE_PASSWORD:
		err = ChangePassword(req, c, gm)
	case requests.C2S_MARKETPLACE_JOBREFILL:
		err = MarketplaceJobRefill(req, c, gm)
	case requests.C2S_COOP_START:
		err = CoopStart(req, c, gm)
	case requests.C2S_COOP_ACTIVELIST:
		err = CoopActiveList(req, c, gm)
	case requests.C2S_FASTFOOD_COOK:
		err = FastFoodCook(req, c, gm)
	case requests.C2S_FASTFOOD_NPC:
		err = FastFoodCook(req, c, gm)
	case requests.C2S_CHANGE_AVATAR:
		err = ChangeAvatar(req, c, gm)
	case requests.C2S_CAFE_RECOOK:
		err = Recook(req, c, gm)
	case requests.C2S_CAFE_TUTORIAL_FINISH:
		err = TutorialComplete(c, gm)
	case requests.C2S_SPECIAL_EVENT:
		err = SendSpecialEvent(c, gm)
	case requests.C2S_SHOP_CARRIER_PIGEON:
		err = BuyIngredientFromShopCarrier(req, c, gm)
	case requests.C2S_HIGHSCORE_LIST:
		err = SendHighscoreList(req, c, gm)
	case requests.C2S_COOP_DETAIL:
		err = CoopDetail(req, c, gm)
	case requests.C2S_COOP_LEAVE:
		err = CoopLeave(req, c, gm)
	case requests.C2S_COOP_JOIN:
		err = CoopJoin(req, c, gm)
	case requests.C2S_COOP_EXTEND:
		err = CoopExtend(req, c, gm)
	case requests.C2S_MARKETPLACE_SEEKINGJOB:
		err = SeekingJob(req, c, gm)
	default:
		log.Infof("NOT IMPLEMENTED: %v", req.Args[0])
	}

	if err != nil {
		return fmt.Errorf("Error during command handling: %s", err)
	}

	return nil
}
