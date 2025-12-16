package commands

import (
	"cafego/internal/client"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	RegisterCommand(requests.C2S_MARKETPLACE_SEEKINGJOB,
		CommandConfig{
			Name:       "SeekingJob",
			Identifier: responses.S2C_MARKETPLACE_SEEKINGJOB,
			MinArgs:    3,
			MaxArgs:    3,
			IsBool:     true,
		},
		SeekingJobValidator,
		SeekingJob,
	)
}

func SeekingJob(req *requests.Request, c *client.Client, gm *managers.GameManager) error {
	isSeekingJob := utils.If(req.Args[2] == "1", true, false)
	if !isSeekingJob {
		c.Player.SeekingJob = false
		c.Player.ClearOffers()
	} else {
		c.Player.SeekingJob = true
	}

	c.Location.Broadcast("mts", "-1", "0", strconv.Itoa(c.Player.ID), utils.If(c.Player.SeekingJob, "1", "0"))

	return nil

}

func SeekingJobValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm CommandConfig) (string, ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), MAX_ARGS
		}
	}

	if cm.IsBool {
		if req.Args[2] != "0" && req.Args[2] != "1" {
			return fmt.Sprintf("Invalid args for boolean: %v", req.Args[2]), INVALID_ARGS
		}
	}

	isSeekingJob := utils.If(req.Args[2] == "1", true, false)
	if !isSeekingJob {
		if c.Player.OpenJobs <= 0 {
			return "Cant seek job, because player is out of jobs", NOT_DECLARED
		}
	}
	return "Command ran without any errors.", NO_ERROR
}
