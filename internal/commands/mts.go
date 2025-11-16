package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"strconv"
)

func SeekingJob(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	seekingStatus, err := strconv.Atoi(req.Args[2])
	if err != nil {
		return err
	}

	if seekingStatus != 1 {
		c.Player.SeekingJob = false
	} else {
		if c.Player.OpenJobs > 0 {
			c.Player.SeekingJob = true
		} else {
			// no errors defined when the player is out open jobs
			return nil
		}
	}

	c.Location.Broadcast("mts", "-1", "0", strconv.Itoa(c.Player.ID), strconv.Itoa(seekingStatus))

	return nil

}
