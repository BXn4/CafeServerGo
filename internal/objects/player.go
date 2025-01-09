package objects

import (
  "strings"
  "strconv"
  "cafego/internal/utils"
)

type Player struct {
	ID                  int
	Cash                int
	Gold                int
	XP                  int
	InstantCookings     int
	OpenJobs            int
	PlayedWheel         bool
	AllowFriendRequests bool
	AllowEmails         bool
	EmailVerified       bool
	NewGifts            int
	Username            string
	Avatar              Avatar
	Position            []int
	Mastery             string // TODO: Create proper mastery
  WorkTimeLeft        int
  SeekingJob          bool
}

func (player *Player) String() string{
    params := []string{
      strconv.Itoa(player.ID),
      strconv.Itoa(player.ID),
      strconv.Itoa(player.XP),
      strconv.Itoa(player.Position[0]),
      strconv.Itoa(player.Position[1]),
      strconv.Itoa(player.WorkTimeLeft),
      strconv.Itoa(player.OpenJobs),
      utils.If(player.SeekingJob, "1", "0"),
      utils.If(player.AllowFriendRequests, "1", "0"),
      player.Avatar.String(),
    }
  return strings.Join(params, "+")
}



