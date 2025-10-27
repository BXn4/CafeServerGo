package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/models/player"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

// lgn - Login
func Login(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	name := req.Args[2]
	password := req.Args[3]

	// Check credentials
	p, statusCode, err := c.DB.Authenticate(name, password)

	// Check if already logged in disconnect the client
	if err == nil {
		if searched, _ := gm.GetClientByName(name); searched != nil {
			statusCode = 15
			// Check if the same ip because than its most likely a bug
			/* if c.GetIP() == searched.GetIP() {
			err = searched.Disconnect() // Kick client out
			} */
		}
	}

	statusCodeStr := strconv.Itoa(statusCode)

	// Send login response (lgn)
	c.SendExtensionResponse("lgn", "1", statusCodeStr)
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

		// Send room list (rlu)
		err = RoomList(req, c, gm)
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

		// Handle login bonus (lbu)
		err = LoginRewards(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nlbu request: %v", err)
		}

		// Send Ping (pin)
		err = SendPing(req, c, gm)
		if err != nil {
			return fmt.Errorf("\npin request: %v", err)
		}

		// Send Friends
		err = SendFriendsAvatar(req, c, gm)
		if err != nil {
			return fmt.Errorf("\nbga request: %v", err)
		}

		p.OnAchievementEarned = func(id int, level int, p *player.Player) {
			gm.SendEarnAchievement(id, level, p.Username)
		}

		p.MakeAchievementCurrentLevels() // only make if its not exist

		c.Player.IsTutorialCompleted = true // Default false, because after register, the customers should not start.
		// And we cant trigger the tutorial event in the game, so its just works after register. If the player disconnects after the tutorial, we cant trigger it.
	}

	return nil
}
