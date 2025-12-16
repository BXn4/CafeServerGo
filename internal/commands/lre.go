package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func init() {
	RegisterCommand(requests.C2S_REGISTER,
		CommandConfig{
			Name:       "Register",
			Identifier: responses.S2C_REGISTER,
			MinArgs:    19,
			MaxArgs:    19,
		},
		RegisterValidator,
		Register,
	)
}

// TODO: add more chars
var invalidChars = "+%&*/()[]{}\"'\\´`^°§€²³.,;:?µ$"

func Register(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	username := req.Args[2]
	email := req.Args[3]
	password := req.Args[4]

	log.Debug("Everything is fine! Player register should start")

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

	player.IsRegistered = true
	player.DailyLogin = time.Now().UTC()
	c.DB.SetRegistered(player.ID)
	c.DB.UpdateDailyLogin(player.ID, player.DailyLogin)
	return nil
}

func RegisterValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	username := req.Args[2]
	email := req.Args[3]
	password := req.Args[4]
	acceptedTerms := req.Args[5]

	// TODO
	log.Warn("SOME REGISTER COMMANDS NOT USED IN THE GAME, BUT EXIST IN THE LANGUAGE FILE! USERNAME SHORT ETC. TODO!")

	// Check if username contains unique characters
	if username == "." || strings.ContainsAny(username, invalidChars) {
		return "Can't register the player, because the username is wrong / contains not allowed chars!", USERNAME_WRONG
	}

	// Check if username is short
	if len(username) < 4 {
		return "Can't register the player, because the username is short!", USERNAME_SHORT
	}
	// Check if username is long
	if len(username) > 24 {
		return "Can't register the player, because the username is long!", USERNAME_LONG
	}

	// TODO: Check if username contains a bad word
	/* if false {
	c.SendExtensionResponse("lre", "-1", BAD_WORD)
	return nil
	} */

	// Check if email is valid
	if email == "." || strings.ContainsAny(strings.Split(email, "@")[0], invalidChars) {
		return "Can't register the player, because the email format is invalid!", EMAIL_INVALID
	}

	// ACCEPT_TERMS
	if acceptedTerms == "0" {
		return "Can't register the player, because the player not accepted the terms!", ACCEPT_TERMS
	}

	// PASSWORD_SHORT
	if len(password) < 4 {
		return "Can't register the player, because the password is short!", PASSWORD_SHORT
	}

	// PASSWORD_INVALID
	if strings.ContainsAny(password, invalidChars) {
		return "Can't register the player, because the password is containts invalid chars!", PASSWORD_INVALID
	}

	// EMAIL_LONG
	if len(email) > 320 {
		return "Can't register the player, because the email is long!", EMAIL_LONG

	}
	// EMAIL_INVALID
	_, err := mail.ParseAddress(email)
	if err != nil {
		return "Can't register the player, because the email is invalid", EMAIL_INVALID
	}

	// ACCOUNT_EXIST
	_, err = c.DB.GetPlayerByName(username)
	if err == nil {
		return "Can't register the player, because the username, or email exist!", ACCOUNT_EXIST
	}

	return "Command ran without any errors.", NO_ERROR
}
