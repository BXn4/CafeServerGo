package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"fmt"
	"strconv"
)

// lgn - Login
func Login(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	name := req.Args[2]
	password := req.Args[3]

	// Check credentials
	statusCode, err := c.DB.Authenticate(name, password)

	// Check if already logged in disconnect the client
	if err == nil {
		if searched, _ := gm.GetClientByName(name); searched != nil {
			statusCode = 15
			// Check if the same ip because than its most likely a bug
			if c.GetIP() == searched.GetIP() {
				err = searched.Disconnect() // Kick client out
				println("STCUK HERE? 1")
			}
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

	// Send room list (rlu)
	err = RoomList(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\trlu request: %v", err)
	}

	// Send user info (gui)
	err = UserInfo(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tgui request: %v", err)
	}

	// Send balancing constants (sbc)
	err = SendBalancingConstant(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tsbc request: %v", err)
	}

	// Send mastery info (lmi)
	err = SendMasteryInfo(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tlmi request: %v", err)
	}

	// Send fridge info (ifr)
	err = SendFridgeInventory(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tifr request: %v", err)
	}

	// Handle login bonus (lbu)
	err = LoginRewards(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tlbu request: %v", err)
	}

	// Send Ping (pin)
	err = SendPing(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tpin request: %v", err)
	}

	// Send Friends
	err = SendFriendsAvatar(req, c, gm)
	if err != nil {
		return fmt.Errorf("\n\tbga request: %v", err)
	}

	return nil
}
