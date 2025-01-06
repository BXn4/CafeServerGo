package responses

import (
  "strings"
)

// Simple static responses
const (
  // System requests
  POLICY_FILE string = "<cross-domain-policy><allow-access-from domain='*' to-ports='*' /></cross-domain-policy>"
  VERSION_CHECK = "<msg t='sys'><body action='apiOK' r='0'></body></msg>"
  AUTO_JOIN = "<msg t='sys'><body action='joinOK' r='1'><pid id='0'/><vars /><uLs r='1'><u i='2' m='0'><n><![CDATA[]]></n><vars></vars></u></uLs></body></msg>"
  ROUND_TRIP = "<msg t='sys'><body action='roundTripRes' r='1'></body></msg>"
  LOGOUT = "<msg t='sys'><body action='logout' r='0'></body></msg>"

  // Extension responses
  S2C_CREATE_AVATAR = "lca"
  S2C_REGISTER = "lre"
  S2C_LOGIN = "lgn"
  S2C_PLAY_WITHOUT_REGISTER = "lwr"
  S2C_COMEBACK_BONUS = "lcb"
  S2C_CHANGE_AVATAR = "cha"
  S2C_LOST_PASSWORD = "llp"
  S2C_CHANGE_PASSWORD = "lcp"
  S2C_CAFE = "sgc"
  S2C_USER_INFO = "gui"
  S2C_INVENTORY_FRIDGE = "ifr"
  S2C_ASSETS_SYNCHRONIZE = "asy"
  S2C_SERVER_MESSAGE = "sms"
  S2C_CAFE_COOK = "ccc"
  S2C_CAFE_INSTANTCOOK = "cic"
  S2C_CAFE_CLEAN = "ccn"
  S2C_CAFE_STOVE_DELIVER_INFO = "csi"
  S2C_CAFE_STOVE_DELIVER = "csd"
  S2C_CAFE_WALK = "cwa"
  S2C_CAFE_CHAT = "cch"
  S2C_MARKETPLACE_JOIN = "mjm"
  S2C_MARKETPLACE_JOB = "mjo"
  S2C_MARKETPLACE_JOBREFILL = "mjr"
  S2C_MARKETPLACE_SEEKINGJOB = "mts"
  S2C_JOB_PAYCHECK = "wpc"
  S2C_JOB_USER_ACTION = "wua"
  S2C_SHOP_AVAILIBILITY = "sga"
  S2C_SHOP_BUY_ITEM = "sbi"
  S2C_SHOP_DELETE_ITEM = "sdi"
  S2C_EDITOR_INVENTORY = "ein"
  S2C_EDITOR_MODE = "edi"
  S2C_EDITOR_MOVE_OBJECT = "emo"
  S2C_EDITOR_ROTATE_OBJECT = "ero"
  S2C_EDITOR_EXTEND = "eex"
  S2C_EDITOR_BUY_OBJECT = "ebu"
  S2C_EDITOR_STORE_OBJECT = "est"
  S2C_EDITOR_SELL_OBJECT = "ese"
  S2C_EDITOR_BUY_FLOOR = "ebf"
  S2C_NPC_AVATAR = "nav"
  S2C_NPC_HIRE = "nhi"
  S2C_NPC_FIRE = "nfi"
  S2C_NPC_CUSTOMIZE = "ncu"
  S2C_NPC_ACTION = "nac"
  S2C_JOIN_USERLIST = "jul"
  S2C_JOIN_USERJOIN = "juj"
  S2C_JOIN_USERQUIT = "juq"
  S2C_OTHERPLAYER_INFO = "bop"
  S2C_INVITE_FRIEND = "cif"
  S2C_SHOP_CARRIER_PIGEON = "scp"
  S2C_CASH_HASH = "gch"
  S2C_LOGIN_BONUS = "lbu"
  S2C_BUDDY_AVATARS = "bga"
  S2C_BUDDY_INGAME = "big"
  S2C_SOCIAL_BUDDIES = "sbs"
  S2C_SEND_BLANCING_CONSTANTS = "sbc"
  S2C_CAFE_ACHIEVEMENT_LIST = "cal"
  S2C_CAFE_ACHIEVEMENT_EARN = "cae"
  S2C_LOGIN_FEATURES = "lfe"
  S2C_JOIN_CAFE = "jca"
  S2C_VERSION_CHECK = "vck"
  S2C_KICK_USER = "cku"
  S2C_CHAT_PUNISHMENT = "cpu"
  S2C_GIFT_REMOVE = "grm"
  S2C_GIFT_USE = "gus"
  S2C_GIFT_PLAYERGIFTS = "gmg"
  S2C_GIFT_SENDABLEGIFTS = "gag"
  S2C_GIFT_ALLREADYSEND_PLAYERS = "gap"
  S2C_COOP_DETAIL = "cod"
  S2C_COOP_ACTIVELIST = "coa"
  S2C_COOP_START = "cos"
  S2C_COOP_JOIN = "coj"
  S2C_COOP_LEAVE = "col"
  S2C_COOP_EXTEND = "cox"
  S2C_COOP_FINISH = "cof"
  S2C_HIGHSCORE_LIST = "hsl"
  S2C_ALLOW_BUDDY_REQUESTS = "abr"
  S2C_ALLOW_MAIL_REQUESTS = "amr"
  S2C_SPECIAL_EVENT = "see"
  S2C_SOCIAL_LOGIN_BONUS = "slb"
  S2C_SOCIAL_TRIGGEREVENT = "ste"
  S2C_CAFE_RECOOK = "crc"
  S2C_MASTERY_INFO = "lmi"
  S2C_SPECIAL_OFFER_EVENT = "soe"
  S2C_PAYMENT_SHOP_PRICE_CHANGE = "ppc"
  S2C_MINI_MUFFIN = "mmu"
  S2C_MINI_MUFFIN_GUEST = "mmg"
  S2C_WHEELOFFORTUNE = "mwf"
  S2C_FASTFOOD_COOK = "ffc"
  S2C_FASTFOOD_NPC = "ffn"
  S2C_EMAIL_CONFIRMED = "emc"
  S2C_EMAIL_VERIFICATION = "emv"

  S2C_ROOMLIST = "rlu"
  S2C_JOIN_ROOM = "jro"
  S2C_PING = "pin"
)

type Response struct {
  Args []string
}
//TODO: This might need to change

func WrapSystemResponse(args ...string) string {
  return strings.Join(args, "") + "\x00"
}

func WrapExtensionResponse(args ...string) string {
  return "%xt%" + strings.Join(args, "%") + "%\x00"
}


