package requests

import (
	"errors"
	"strings"
)

type RequestKind int

const (
	// System request types
	UNKNOWN RequestKind = iota
	POLICY_FILE
	DISCONNECT
	VERSION_CHECK
	AUTO_JOIN
	ROUND_TRIP
	LOGIN

	// Extension request types
	C2S_LOGIN
	C2S_LOGIN_SOCIAL
	C2S_REGISTER
	C2S_VERSION_CHECK
	C2S_CREATE_AVATAR
	C2S_LOST_PASSWORD
	C2S_CHANGE_PASSWORD
	C2S_PLAY_WITHOUT_REGISTER
	C2S_CHANGE_AVATAR
	C2S_CAFE_COOK
	C2S_CAFE_INSTANTCOOK
	C2S_CAFE_CLEAN
	C2S_CAFE_STOVE_DELIVER_INFO
	C2S_CAFE_STOVE_DELIVER
	C2S_CAFE_WALK
	C2S_CAFE_CHAT
	C2S_MARKETPLACE_JOIN
	C2S_MARKETPLACE_JOB
	C2S_JOB_PAYCHECK
	C2S_JOB_USER_ACTION
	C2S_MARKETPLACE_JOBREFILL
	C2S_MARKETPLACE_SEEKINGJOB
	C2S_OTHERPLAYER_INFO
	C2S_LOGIN_FEATURES
	C2S_SHOP_AVAILIBILITY
	C2S_SHOP_BUY_ITEM
	C2S_SHOP_DELETE_ITEM
	C2S_EDITOR_MODE
	C2S_EDITOR_MOVE_OBJECT
	C2S_EDITOR_ROTATE_OBJECT
	C2S_EDITOR_EXTEND
	C2S_EDITOR_BUY_OBJECT
	C2S_EDITOR_STORE_OBJECT
	C2S_EDITOR_SELL_OBJECT
	C2S_EDITOR_BUY_FLOOR
	C2S_NPC_HIRE
	C2S_NPC_FIRE
	C2S_NPC_CUSTOMIZE
	C2S_INVITE_FRIEND
	C2S_SHOP_CARRIER_PIGEON
	C2S_CASH_HASH
	C2S_SOCIAL_BUDDIES
	C2S_BUDDY_INGAME
	C2S_GIFT_SEND
	C2S_GIFT_REMOVE
	C2S_GIFT_USE
	C2S_GIFT_PLAYERGIFTS
	C2S_GIFT_SENDABLEGIFTS
	C2S_GIFT_ALLREADYSEND_PLAYERS
	C2S_CAFE_ACHIEVEMENT_LIST
	C2S_CHANGE_LANGUAGE
	C2S_CAFE_TUTORIAL_FINISH
	C2S_JOIN_CAFE
	C2S_KICK_USER
	C2S_REPORT_PLAYER
	C2S_ALLOW_BUDDY_REQUESTS
	C2S_ALLOW_MAIL_REQUESTS
	C2S_COOP_DETAIL
	C2S_COOP_ACTIVELIST
	C2S_COOP_START
	C2S_COOP_JOIN
	C2S_COOP_LEAVE
	C2S_COOP_EXTEND
	C2S_HIGHSCORE_LIST
	C2S_SPECIAL_EVENT
	C2S_CAFE_RECOOK
	C2S_MINI_MUFFIN
	C2S_WHEELOFFORTUNE
	C2S_FASTFOOD_COOK
	C2S_FASTFOOD_NPC
	C2S_EMAIL_VERIFICATION
	C2S_ROOMLIST
	C2S_PING
)

func LookupRequestKind(kindStr string) RequestKind {
	kindLookup := map[string]RequestKind{
		// System
		"logout":    DISCONNECT,
		"verChk":    VERSION_CHECK,
		"autoJoin":  AUTO_JOIN,
		"roundTrip": ROUND_TRIP,
		"login":     LOGIN,

		// Extension
		"lgn": C2S_LOGIN,
		"lgs": C2S_LOGIN_SOCIAL,
		"lre": C2S_REGISTER,
		"vck": C2S_VERSION_CHECK,
		"lca": C2S_CREATE_AVATAR,
		"llp": C2S_LOST_PASSWORD,
		"lcp": C2S_CHANGE_PASSWORD,
		"lwr": C2S_PLAY_WITHOUT_REGISTER,
		"cha": C2S_CHANGE_AVATAR,
		"ccc": C2S_CAFE_COOK,
		"cic": C2S_CAFE_INSTANTCOOK,
		"ccn": C2S_CAFE_CLEAN,
		"csi": C2S_CAFE_STOVE_DELIVER_INFO,
		"csd": C2S_CAFE_STOVE_DELIVER,
		"cwa": C2S_CAFE_WALK,
		"cch": C2S_CAFE_CHAT,
		"mjm": C2S_MARKETPLACE_JOIN,
		"mjo": C2S_MARKETPLACE_JOB,
		"wpc": C2S_JOB_PAYCHECK,
		"wua": C2S_JOB_USER_ACTION,
		"mjr": C2S_MARKETPLACE_JOBREFILL,
		"mts": C2S_MARKETPLACE_SEEKINGJOB,
		"bop": C2S_OTHERPLAYER_INFO,
		"lfe": C2S_LOGIN_FEATURES,
		"sga": C2S_SHOP_AVAILIBILITY,
		"sbi": C2S_SHOP_BUY_ITEM,
		"sdi": C2S_SHOP_DELETE_ITEM,
		"edi": C2S_EDITOR_MODE,
		"emo": C2S_EDITOR_MOVE_OBJECT,
		"ero": C2S_EDITOR_ROTATE_OBJECT,
		"eex": C2S_EDITOR_EXTEND,
		"ebu": C2S_EDITOR_BUY_OBJECT,
		"est": C2S_EDITOR_STORE_OBJECT,
		"ese": C2S_EDITOR_SELL_OBJECT,
		"ebf": C2S_EDITOR_BUY_FLOOR,
		"nhi": C2S_NPC_HIRE,
		"nfi": C2S_NPC_FIRE,
		"ncu": C2S_NPC_CUSTOMIZE,
		"cif": C2S_INVITE_FRIEND,
		"scp": C2S_SHOP_CARRIER_PIGEON,
		"gch": C2S_CASH_HASH,
		"sbs": C2S_SOCIAL_BUDDIES,
		"big": C2S_BUDDY_INGAME,
		"gse": C2S_GIFT_SEND,
		"grm": C2S_GIFT_REMOVE,
		"gus": C2S_GIFT_USE,
		"gmg": C2S_GIFT_PLAYERGIFTS,
		"gag": C2S_GIFT_SENDABLEGIFTS,
		"gap": C2S_GIFT_ALLREADYSEND_PLAYERS,
		"cal": C2S_CAFE_ACHIEVEMENT_LIST,
		"clg": C2S_CHANGE_LANGUAGE,
		"ctf": C2S_CAFE_TUTORIAL_FINISH,
		"jca": C2S_JOIN_CAFE,
		"cku": C2S_KICK_USER,
		"rpl": C2S_REPORT_PLAYER,
		"abr": C2S_ALLOW_BUDDY_REQUESTS,
		"amr": C2S_ALLOW_MAIL_REQUESTS,
		"cod": C2S_COOP_DETAIL,
		"coa": C2S_COOP_ACTIVELIST,
		"cos": C2S_COOP_START,
		"coj": C2S_COOP_JOIN,
		"col": C2S_COOP_LEAVE,
		"cox": C2S_COOP_EXTEND,
		"hsl": C2S_HIGHSCORE_LIST,
		"see": C2S_SPECIAL_EVENT,
		"crc": C2S_CAFE_RECOOK,
		"mmu": C2S_MINI_MUFFIN,
		"mwf": C2S_WHEELOFFORTUNE,
		"ffc": C2S_FASTFOOD_COOK,
		"ffn": C2S_FASTFOOD_NPC,
		"emv": C2S_EMAIL_VERIFICATION,
		"rlu": C2S_ROOMLIST,
		"pin": C2S_PING,
	}

	if val, ok := kindLookup[kindStr]; ok {
		return val
	}

	return UNKNOWN
}

type Request struct {
	Kind RequestKind
	Args []string
}

func ParseRequest(raw_request string) (*Request, error) {
	trimmed_req := strings.TrimSpace(raw_request)

	// Catch if policy-file-request
	if strings.HasPrefix(trimmed_req, "<policy-file-request/>") {
		return &Request{
			Kind: POLICY_FILE,
		}, nil
	}

	// Check if sys request or extension
	if strings.HasPrefix(trimmed_req, "<msg t='sys'>") {
		return ParseSystemRequest(trimmed_req)
	} else if strings.HasPrefix(trimmed_req, "%xt") {
		return ParseExtensionRequest(trimmed_req)
	}

	return nil, errors.New("Can not parse request!")
}

func ParseSystemRequest(req string) (*Request, error) {
	rawParams := strings.Split(strings.Split(req, "<body ")[1], ">")[0]
	paramPairs := strings.Split(rawParams, " ")
	var params []string
	var action string
	for _, paramPair := range paramPairs {
		// Save params
		trimmed_pair := strings.ReplaceAll(paramPair, "'", "")
		params = append(params, trimmed_pair)

		// Save action
		if strings.HasPrefix(trimmed_pair, "action") {
			action = strings.Split(trimmed_pair, "=")[1]
		}
	}

	// Look up request type
	kind := LookupRequestKind(action)

	return &Request{
		Kind: kind,
		Args: params,
	}, nil
}

func ParseExtensionRequest(req string) (*Request, error) {

	trimmed_req := strings.Trim(req, "%")
	args := strings.Split(trimmed_req, "%")
	args = args[2:] // Cut of the 'xt' and the 'CafeEx'

	kind := LookupRequestKind(args[0])

	return &Request{
		Kind: kind,
		Args: args,
	}, nil
}
