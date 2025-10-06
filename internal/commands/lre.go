package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"net/mail"
	"strings"

	"github.com/charmbracelet/log"
)

const (
	USERNAME_WRONG   = "1"
	USERNAME_SHORT   = "2"
	USERNAME_LONG    = "3"
	EMAIL_WRONG      = "4"
	PASSWORD_WRONG   = "5"
	ACCEPT_TERMS     = "98"
	PASSWORD_SHORT   = "96"
	PASSWORD_INVALID = "10"
	EMAIL_LONG       = "94"
	EMAIL_INVALID    = "14"
	ACCOUNT_EXIST    = "13"
	BAD_WORD         = "93"
)

// TODO: add more chars
var invalidChars = "+%&*/()[]{}\"'\\´`^°§€²³.,;:?µ$"

func Register(req *requests.Request, c *client.Client, gm *managers.GameManager) error {

	username := req.Args[2]
	email := req.Args[3]
	password := req.Args[4]
	acceptedTerms := req.Args[5]

	// Check if username contains unique characters
	if username == "." || strings.ContainsAny(username, invalidChars) {
		c.SendExtensionResponse("lre", "-1", USERNAME_WRONG)
		return nil
	}

	// Check if username is short
	if len(username) < 4 {
		c.SendExtensionResponse("lre", "-1", USERNAME_SHORT)
		return nil
	}
	// Check if username is long
	if len(username) > 24 {
		c.SendExtensionResponse("lre", "-1", USERNAME_LONG)
		return nil
	}

	// Check if email is valid
	if email == "." || strings.ContainsAny(strings.Split(email, "@")[0], invalidChars) {
		c.SendExtensionResponse("lre", "-1", EMAIL_WRONG)
		return nil
	}

	// PASSWORD_WRONG
	if password == "." {
		c.SendExtensionResponse("lre", "-1", PASSWORD_WRONG)
		return nil
	}

	// ACCEPT_TERMS
	if acceptedTerms == "0" {
		c.SendExtensionResponse("lre", "-1", ACCEPT_TERMS)
		return nil
	}

	// PASSWORD_SHORT
	if len(password) < 4 {
		c.SendExtensionResponse("lre", "-1", PASSWORD_SHORT)
		return nil
	}

	// PASSWORD_INVALID
	if strings.ContainsAny(password, invalidChars) {
		c.SendExtensionResponse("lre", "-1", PASSWORD_INVALID)
		return nil
	}

	// EMAIL_LONG
	if len(email) > 320 {
		c.SendExtensionResponse("lre", "-1", EMAIL_LONG)
		return nil
	}
	// EMAIL_INVALID
	_, err := mail.ParseAddress(email)
	if err != nil {
		c.SendExtensionResponse("lre", "-1", EMAIL_INVALID)
		return nil
	}

	// ACCOUNT_EXIST
	_, err = c.DB.GetPlayerByName(username)
	if err == nil {
		c.SendExtensionResponse("lre", "-1", ACCOUNT_EXIST)
		return nil
	}

	log.Debug("Everything is fine! Player register should start")

	// TODO: Check if username contains a bad word
	/* if false {
	c.SendExtensionResponse("lre", "-1", BAD_WORD)
	return nil
	} */

	avatar := c.Player.Avatar

	log.Debugf("Player avatar: %s", avatar.Apperance())

	player, err := c.DB.CreateAccount(username, email, password, avatar)
	if err != nil {
		return err
	}
	c.Player = player
	location := gm.AddLocation(c.Player.ID)
	c.Location = location

	c.SendExtensionResponse("lre", "-1", "0")

	// Send room list (rlu)
	err = RoomList(req, c, gm)
	if err != nil {
		return err
	}

	// Send user info (gui)
	err = UserInfo(req, c, gm)
	if err != nil {
		return err
	}

	// Send balancing constants (sbc)
	err = SendBalancingConstant(req, c, gm)
	if err != nil {
		return err
	}

	// Send mastery info (lmi)
	err = SendMasteryInfo(req, c, gm)
	if err != nil {
		return err
	}

	// Send fridge info (ifr)
	err = SendFridgeInventory(req, c, gm)
	if err != nil {
		return err
	}

	// Send Ping (pin)
	err = SendPing(req, c, gm)
	if err != nil {
		return err
	}

	return nil
}
