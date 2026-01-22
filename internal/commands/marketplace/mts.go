package marketplace

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/managers"
	"cafego/internal/types/requests"
	"cafego/internal/types/responses"
	"cafego/internal/utils"
	"fmt"
	"strconv"
)

func init() {
	commands.RegisterCommand(requests.C2S_MARKETPLACE_SEEKINGJOB,
		commands.CommandConfig{
			Name:         "SeekingJob",
			Identifier:   responses.S2C_MARKETPLACE_SEEKINGJOB,
			Description:  "Triggers the seeking job feature",
			Args:         "",
			MinArgs:      3,
			MaxArgs:      3,
			IsBool:       true,
			FeatureLevel: 4,
		},
		SeekingJobValidator,
		SeekingJob,
		nil,
	)
}

func SeekingJob(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) error {
	isSeekingJob := utils.If(req.Args[2] == "1", true, false)
	if !isSeekingJob {
		c.Player.SetIsSeekingJob(false)
		c.Player.ClearOffers()
	} else {
		c.Player.SetIsSeekingJob(true)
	}

	c.Location.Broadcast(cm.Identifier, "-1", "0", strconv.Itoa(c.Player.GetID()), utils.If(c.Player.GetIsSeekingJob(), "1", "0"))

	return nil

}

func SeekingJobValidator(req *requests.Request, c *client.Client, gm *managers.GameManager, cm commands.CommandConfig) (string, commands.ErrorCodes) {
	if len(req.Args) < cm.MinArgs {
		return fmt.Sprintf("Not enough args. NEEDED/GOT: %d/%d", cm.MinArgs, len(req.Args)), commands.MIN_ARGS
	}

	if cm.MinArgs > 0 {
		if len(req.Args) > cm.MaxArgs {
			return fmt.Sprintf("Too much args. NEEDED/GOT: %d/%d", cm.MaxArgs, len(req.Args)), commands.MAX_ARGS
		}
	}

	if cm.IsBool {
		if req.Args[2] != "0" && req.Args[2] != "1" {
			return fmt.Sprintf("Invalid args for boolean: %v", req.Args[2]), commands.INVALID_ARGS
		}
	}

	if c.Player.GetLevel() < cm.FeatureLevel {
		return "Cant join the marketplace, because the player not yet reached the feature.", commands.NOT_DECLARED
	}

	isSeekingJob := utils.If(req.Args[2] == "1", true, false)
	if !isSeekingJob {
		if c.Player.GetOpenJobs() <= 0 {
			return "Cant seek job, because player is out of jobs", commands.NOT_DECLARED
		}
	}
	return "Command ran without any errors.", commands.NO_ERROR
}
