package player

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/commands/cafe"
	"cafego/internal/commands/friends"
	"cafego/internal/managers"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"fmt"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
)

func init() {
	commands.RegisterCommand(requests.C2S_LOGIN,
		commands.CommandConfig{
			Name:        "Login",
			Identifier:  responses.S2C_LOGIN,
			Description: "Login",
			Args:        "{}",
			MinArgs:     2,
			MaxArgs:     2,
		},
		nil,
		Login,
		nil,
	)
}

// lgn - Login
func Login(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	name := req.Args[2]
	password := req.Args[3]

	// Check credentials
	p, statusCode, err := c.DB.Authenticate(name, password)

	// Check if already logged in disconnect the client
	if err == nil {
		log.Infof("Checking if user %s is already logged in", name)
		searched, searchErr := gm.GetClientByName(name)
		if searched != nil && searchErr == nil {
			log.Infof("User %s found as already logged in on client %d, kicking existing session", name, searched.ClientID)
			statusCode = 15

			// Check if same IP and kick the existing user
			if c.GetIP() == searched.GetIP() {
				kickErr := searched.Disconnect()
				if kickErr != nil {
					log.Errorf("Failed to kick existing session for user %s: %v", name, kickErr)
				} else {
					log.Infof("Successfully kicked existing session for user %s", name)
				}
				// Small delay to ensure cleanup completes
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			log.Infof("User %s is not currently logged in (searchErr: %v)", name, searchErr)
		}
	}

	statusCodeStr := strconv.Itoa(statusCode)

	// Send login response (lgn)
	c.SendExtensionResponse(cm.Identifier, "1", statusCodeStr)
	if statusCode != 0 {
		if statusCode == 15 {
			return fmt.Errorf("Player %v is already logged in", name)
		} else {
			return fmt.Errorf("Access denied")
		}
	}

	if p != nil {
		// Set player
		c.Player = p

		p.MakeAchievementCurrentLevels() // only make if its not exist

		// Handle login bonus (lbu)
		err = LoginRewards(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nlbu request: %v", err)
		}

		// Send room list (rlu)
		err = cafe.RoomList(req, c, gm, cm) // -- cm is not used
		if err != nil {
			return fmt.Errorf("\nrlu request: %v", err)
		}

		// Send user info (gui)
		err = UserInfo(req, c, gm)
		if err != nil {
			return fmt.Errorf("\ngui request: %v", err)
		}

		// Send balancing constants (sbc)
		err = SendBalancingConstant(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nsbc request: %v", err)
		}

		// Send mastery info (lmi)
		err = SendMasteryInfo(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nlmi request: %v", err)
		}

		// Send Ping (pin)
		err = SendPing(req, c, gm, cm) // -- cm is not used
		if err != nil {
			return fmt.Errorf("\npin request: %v", err)
		}

		// Send Friends
		err = friends.SendFriendsAvatar(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nbga request: %v", err)
		}

		p.OnAchievementEarned = func(id int, level int, p *player.Player) {
			gm.SendEarnAchievement(id, level, p.GetUsername())
		}

		if c.Player.GetXP() == 0 && !c.Player.GetIsTutorialCompleted() {
			// reset it to allow player to complete the tutorial
			c.SendExtensionResponse("lgn", "-1", strconv.Itoa(commands.LOGIN_SET_FIRST_LOGIN))
		} else {
			c.Player.SetIsTutorialCompleted(true) // Default false, because after register, the customers should not start.
		}
	}

	return nil
}
