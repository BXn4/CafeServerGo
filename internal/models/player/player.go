package player

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/gift"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

type Player struct {
	ID                  int             `gorm:"primaryKey;autoIncrement;type:int"`
	Email               string          `gorm:"not null;type:text"`
	Password            string          `gorm:"not null;type:text"`
	Cash                int             `gorm:"column:cash;type:int"`
	Gold                int             `gorm:"column:gold;type:int"`
	XP                  int             `gorm:"column:xp;type:int"`
	InstantCookings     int             `gorm:"column:instant_cookings;type:int"`
	OpenJobs            int             `gorm:"column:open_jobs;type:int"`
	PlayedWheel         bool            `gorm:"column:played_wheel;default:0;type:bool"`
	AllowFriendRequests bool            `gorm:"column:allow_friend_requests;default:0;type:bool"`
	Friends             simple.IntSlice `gorm:"column:friends;type:text"`
	FriendsWithGifts    simple.IntSlice `gorm:"column:friends_with_gifts;type:text"`
	SendableGifts       gift.GiftList   `gorm:"column:sendable_gifts;type:text"`
	AllowEmails         bool            `gorm:"column:allow_emails;default:0;type:bool"`
	EmailVerified       bool            `gorm:"column:email_verified;default:0;type:bool"`
	Username            string          `gorm:"column:username;not null;type:text"`
	IsBanned            bool            `gorm:"column:is_banned;default:0;type:bool"`
	Avatar              avatar.Avatar   `gorm:"column:avatar;type:text"`
	Position            simple.Position `gorm:"-"`
	Mastery             simple.IntMap   `gorm:"column:mastery;type:text"`
	Achievement         simple.IntMap   `gorm:"column:achievement;type:text"`
	WorkTimeLeft        int             `gorm:"-"`
	CoopID              int             `gorm:"column:coop_id;type:int"`
	SeekingJob          bool            `gorm:"-"`
	LastLogin           time.Time       `gorm:"column:last_login;type:datetime"`
	DailyLogin          time.Time       `gorm:"column:daily_login;type:datetime"`
	GiftRefreshTime     time.Time       `gorm:"column:gift_refresh_time;type:datetime"`
	Gifts               gift.GiftList   `gorm:"column:gifts;type:text"`
	AccessLevel         int             `gorm:"column:access_level;default:0;type:int"`
}

func (player *Player) TableName() string {
	return "player"
}

func (p *Player) String() string {
	p.Avatar.Name = p.Username
	params := []string{
		strconv.Itoa(p.ID),
		strconv.Itoa(p.ID),
		strconv.Itoa(p.GetXP()),
		strconv.Itoa(p.Position.X),
		strconv.Itoa(p.Position.Y),
		strconv.Itoa(p.WorkTimeLeft),
		strconv.Itoa(p.OpenJobs),
		utils.If(p.SeekingJob, "1", "0"),
		utils.If(p.AllowFriendRequests, "1", "0"),
		p.Avatar.String(),
	}
	return strings.Join(params, "+")
}

func (p *Player) GetLevel() int {
	if p.XP < 90 {
		return int(math.Pow(math.Floor(float64(p.XP)/10), 1/2))
	}
	return int(math.Pow(math.Floor(float64(p.XP)/5), 1/3.72))
}

func (p *Player) GetUsername() string {
	return p.Avatar.Name
}

func (p *Player) UpdateMastery(dishID int) {
	p.Mastery[dishID] += 1
}

func (p *Player) GetDishMasteryLevel(dishID int) int {
	// Get items info
	timesCooked := p.Mastery[dishID]
	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		log.Printf("Invalid dish ID: %v", err)
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

func (p *Player) GetDishMasteryServings(dishID int) int {
	masteryLevel := p.GetDishMasteryLevel(dishID)

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		log.Printf("Invalid dish ID: %v", err)
		return 0 // ???
	}

	baseServings := dishInfo.Servings
	if masteryLevel < 1 {
		return baseServings
	} else {
		return int(math.Round(float64(baseServings) * 1.05))
	}
}

func (p *Player) GetDishMasteryXP(dishID int) int {
	masteryLevel := p.GetDishMasteryLevel(dishID)

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		log.Printf("Invalid dish ID: %v", err)
		return 0 // ???
	}

	baseXP := dishInfo.XP
	if masteryLevel < 2 {
		return baseXP
	} else {
		return int(math.Round(float64(baseXP) * 1.05))
	}
}

func (p *Player) GetDishMasteryDuration(dishID int) int {
	masteryLevel := p.GetDishMasteryLevel(dishID)

	dishInfo, err := utils.GetDish(dishID)
	if err != nil {
		log.Printf("Invalid dish ID: %v", err)
		return 0 // ???
	}

	baseDuration := dishInfo.Duration

	if masteryLevel < 3 {
		return baseDuration * 60
	} else {
		return int(math.Round(float64(baseDuration)*0.95) * 60)
	}
}

func (p *Player) AddFriend(id int) {
	p.Friends = append(p.Friends, id)
	p.SetAchivementFriendsCount()
}

func (p *Player) DeleteFriend(id int) {
	index := -1
	for i, f := range p.Friends {
		if f == id {
			index = i
		}
	}
	if index == -1 {
		return
	}
	p.Friends = append(p.Friends[:index], p.Friends[index+1:]...)

}

func (p *Player) AddCash(amount int) {
	p.Cash += amount
	if amount > 0 {
		p.UpdateAchivementSpendChips(amount)
	} else if amount < 0 {
		p.UpdateAchivementSpendChips(amount)
	}
}

func (p *Player) SetCash(amount int) {
	p.Cash = amount
}

func (p *Player) SetGold(amount int) {
	p.Gold = amount
}

func (p *Player) GetCash() int {
	return p.Cash
}

func (p *Player) GetGold() int {
	return p.Gold
}

func (p *Player) AddGold(amount int) {
	p.Gold += amount
	if amount < 0 {
		p.UpdateAchivementSpendGold(amount)
	}
}

func (p *Player) GetXP() int {
	return p.XP
}

func (p *Player) AddXP(amount int) {
	beforeLvl := p.GetLevel()
	p.XP += amount
	if beforeLvl < p.GetLevel() {
		p.AddCash(500)
		p.AddGold(1)
	}
}

func (p *Player) SetXP(amount int) {
	p.XP = amount
}
