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

	// Check if already logged in log him/her out
	if err == nil {
		if searched, _ := gm.GetClientByName(name); searched != nil {
			statusCode = 15
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

	// Handle login bonus (lbu)
	err = LoginRewards(req, c, gm)
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
