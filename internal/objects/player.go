package objects

import (
	"cafego/internal/utils"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	Mastery             map[int]int
	WorkTimeLeft        int
	SeekingJob          bool
}

func (player *Player) String() string {
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

func (p *Player) ParseMastery(mastery string) {
	pairs := strings.Split(mastery, "#")

	p.Mastery = make(map[int]int)

	for _, pair := range pairs {
		parts := strings.Split(pair, "+")

		dishID, err := strconv.Atoi(parts[0])
		if err != nil {
			fmt.Printf("Cant parse mastery dishID to int: %v", err)
		}
		timesCooked, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Printf("Cant parse mastery timesCooked to int: %v", err)
		}
		p.Mastery[dishID] = timesCooked
	}
}

func (p *Player) BuildMastery() string {
	var pairs []string

	for dishID, timesCooked := range p.Mastery {
		pair := fmt.Sprintf("%d+%d", dishID, timesCooked)
		pairs = append(pairs, pair)
	}

	return strings.Join(pairs, "#")
}

func (p *Player) GetDishMasteryLevel(dishID int) int {
	// Get items info
	timesCooked := p.Mastery[dishID]
	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		fmt.Printf("Invalid dish ID: %v", err)
	}

	// Get base duration in minutes
	baseDuration := dishInfo.Duration

	// Get base duration in hours
	durationInHours := float64(baseDuration) / 60.0

	// TODO: Comment this shit out WTF ???
	_loc4_ := math.Min(24, float64(durationInHours)*math.Ceil(0.5/float64(durationInHours))*3)

	level1Req := _loc4_ / float64(durationInHours) * 4
	level2Req := _loc4_ / float64(durationInHours) * 20
	level3Req := _loc4_ / float64(durationInHours) * 52

	if timesCooked < int(level1Req) {
		return 0
	} else if timesCooked < int(level2Req) {
		return 1
	} else if timesCooked < int(level3Req) {
		return 2
	} else {
		return 3
	}
}

func (p *Player) GetDishMasteryDuration(dishID int) int {
	masteryLevel := p.GetDishMasteryLevel(dishID)

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		fmt.Printf("Invalid dish ID: %v", err)
	}

	baseDuration := dishInfo.Duration

	if masteryLevel < 3 {
		return baseDuration * 60
	} else {
		return int(math.Floor(float64(baseDuration)*0.95) * 60)
	}
}
