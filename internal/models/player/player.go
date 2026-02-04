package player

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/gift"
	"cafego/internal/models/simple"
	"cafego/internal/utils"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Player struct {
	ID                  int                                           `gorm:"primaryKey;autoIncrement;type:int"`
	Email               string                                        `gorm:"not null;type:text"`
	Password            string                                        `gorm:"not null;type:text"`
	Cash                int                                           `gorm:"column:cash;type:int;default:2000"`
	Gold                int                                           `gorm:"column:gold;type:int;default:11"`
	XP                  int                                           `gorm:"column:xp;type:int;default:0"`
	InstantCookingsUsed int                                           `gorm:"column:instant_cookings_used;type:int;default:0"`
	RefilledJobs        int                                           `gorm:"column:refilled_jobs;type:int;default:0"`
	OpenJobs            int                                           `gorm:"column:open_jobs;type:int;default:0"`
	PlayedWheel         bool                                          `gorm:"column:played_wheel;type:bool;default:false"`
	AllowFriendRequests bool                                          `gorm:"column:allow_friend_requests;stype:bool;default:true"`
	Friends             simple.IntSlice                               `gorm:"column:friends;type:text;default:null"`
	FriendsWithGifts    simple.IntSlice                               `gorm:"column:friends_with_gifts;type:text;default:null"`
	SendableGifts       gift.GiftList                                 `gorm:"column:sendable_gifts;type:text;default:null"`
	AllowEmails         bool                                          `gorm:"column:allow_emails;type:bool;default:true"`
	EmailVerified       bool                                          `gorm:"column:email_verified;type:bool;default:false"`
	Username            string                                        `gorm:"column:username;not null;type:text"`
	IsBanned            bool                                          `gorm:"column:is_banned;default:0;type:bool"`
	Avatar              avatar.Avatar                                 `gorm:"column:avatar;type:text"`
	AvatarChanged       bool                                          `gorm:"column:avatar_changed;type:bool;default:false"`
	position            simple.Position                               `gorm:"-"`
	Mastery             simple.IntMap                                 `gorm:"column:mastery;type:text;default:null"`
	Achievements        simple.IntMap                                 `gorm:"column:achievements;type:text;default:null"`
	achievementsLevel   simple.IntMap                                 `gorm:"-"`
	workTimeLeft        int                                           `gorm:"-"`
	CoopID              int                                           `gorm:"column:coop_id;type:int;default:null"`
	IsStartedCoop       bool                                          `gorm:"column:is_started_coop;type:int;default:false"`
	isSeekingJob        bool                                          `gorm:"-"`
	LastLogin           time.Time                                     `gorm:"column:last_login;type:datetime;default:null"`
	DailyLogin          time.Time                                     `gorm:"column:daily_login;type:datetime;default:null"`
	GiftRefreshTime     time.Time                                     `gorm:"column:gift_refresh_time;type:datetime;default:null"`
	Gifts               gift.GiftList                                 `gorm:"column:gifts;type:text;default:null"`
	IsRegistered        bool                                          `gorm:"column:is_registered;type:bool;default:false"`
	AccessLevel         int                                           `gorm:"column:access_level;default:0;type:int"`
	maxInstants         int                                           `gorm:"-"`
	job                 PlayerJob                                     `gorm:"-"`
	mutex               sync.Mutex                                    `gorm:"-"`
	OnAchievementEarned func(achievementID int, level int, p *Player) `gorm:"-"`
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
		strconv.Itoa(p.position.X),
		strconv.Itoa(p.position.Y),
		strconv.Itoa(p.GetWorkTimeLeft()),
		strconv.Itoa(p.OpenJobs),
		utils.If(p.isSeekingJob, "1", "0"),
		utils.If(p.AllowFriendRequests, "1", "0"),
		p.Avatar.String(),
	}
	return strings.Join(params, "+")
}

func (p *Player) GetIsDailyLogin() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	timePassed := time.Now().UTC().Sub(p.DailyLogin)
	return timePassed >= 24*time.Hour
}

// ** SETTERS ** // ** SETTERS ** // ** SETTERS ** //
func (p *Player) SetID(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.ID = v
}

func (p *Player) SetEmail(v string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Email = v
}

func (p *Player) SetPassword(v string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Password = v
}

func (p *Player) SetCash(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Cash = v
}

func (p *Player) AddCash(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if v > 0 {
		p.UpdateAchivementSpendChips(v)
	} else if v < 0 {
		p.UpdateAchivementSpendChips(v)
	}

	p.Cash += v
}

func (p *Player) SetGold(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Gold = v
}

func (p *Player) AddGold(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	println("adding")
	println(v)
	println("total")
	println(p.Gold)

	p.Gold += v
	if v < 0 {
		p.UpdateAchivementSpendGold(-v)
	}
}

func (p *Player) SetInstantCookings(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.InstantCookingsUsed = v
}

func (p *Player) SetRefilledJobs(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.RefilledJobs = v
}

func (p *Player) AddRefilledJobs() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.RefilledJobs++
}

func (p *Player) SetOpenJobs(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.OpenJobs = v
}

func (p *Player) SetPlayedWheel(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.PlayedWheel = v
}

func (p *Player) SetAllowFriendRequests(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.AllowFriendRequests = v
}

func (p *Player) SetFriends(v simple.IntSlice) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Friends = v
}

func (p *Player) SetFriendsWithGifts(v simple.IntSlice) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.FriendsWithGifts = v
}

func (p *Player) SetSendableGifts(v gift.GiftList) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.SendableGifts = v
}

func (p *Player) SetAllowEmails(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.AllowEmails = v
}

func (p *Player) SetEmailVerified(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.EmailVerified = v
}

func (p *Player) SetUsername(v string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Username = v
}

func (p *Player) SetIsBanned(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.IsBanned = v
}

func (p *Player) SetAvatar(v avatar.Avatar) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Avatar = v
}

func (p *Player) SetAvatarChanged(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.AvatarChanged = v
}

func (p *Player) SetPos(v simple.Position) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.position = v
}

func (p *Player) SetMastery(v simple.IntMap) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Mastery = v
}

func (p *Player) UpdateMastery(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Mastery[v] += 1
}

func (p *Player) SetAchievements(v simple.IntMap) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Achievements = v
}

func (p *Player) SetAchievementsLevel(v simple.IntMap) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.achievementsLevel = v
}

func (p *Player) SetCoopID(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.CoopID = v
}

func (p *Player) SetIsStartedCoop(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.IsStartedCoop = v
}

func (p *Player) SetIsSeekingJob(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.isSeekingJob = v
}

func (p *Player) SetLastLogin(v time.Time) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.LastLogin = v
}

func (p *Player) SetDailyLogin(v time.Time) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.DailyLogin = v
}

func (p *Player) SetGiftRefreshTime(v time.Time) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.GiftRefreshTime = v
}

func (p *Player) SetGifts(v gift.GiftList) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Gifts = v
}

func (p *Player) SetIsRegistered(v bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.IsRegistered = v
}

func (p *Player) SetAccessLevel(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.AccessLevel = v
}

func (p *Player) SetMaxInstants(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.maxInstants = v
}

func (p *Player) SetJob(v PlayerJob) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.job = v
}

func (p *Player) RemoveInstantCooking() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.InstantCookingsUsed++ // need to add, because its checks how much was used.
}

func (p *Player) AddFriend(v int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Friends = append(p.Friends, v)
}

func (p *Player) DeleteFriend(id int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

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

// ** GETTERS ** // ** GETTERS ** // ** GETTERS ** //
func (p *Player) GetID() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.ID
}

func (p *Player) GetEmail() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Email
}

func (p *Player) GetPassword() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Password
}

func (p *Player) GetCash() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Cash
}

func (p *Player) GetGold() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Gold
}

func (p *Player) GetInstantCookingsUsed() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.InstantCookingsUsed
}

func (p *Player) GetRefilledJobs() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.RefilledJobs
}

func (p *Player) GetOpenJobs() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.OpenJobs
}

func (p *Player) GetPlayedWheel() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.PlayedWheel
}

func (p *Player) GetAllowFriendRequests() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.AllowFriendRequests
}

func (p *Player) GetFriends() simple.IntSlice {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Friends
}

func (p *Player) GetFriendsWithGifts() simple.IntSlice {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.FriendsWithGifts
}

func (p *Player) GetSendableGifts() gift.GiftList {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.SendableGifts
}

func (p *Player) GetAllowEmails() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.AllowEmails
}

func (p *Player) GetEmailVerified() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.EmailVerified
}

func (p *Player) GetUsername() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Username
}

func (p *Player) GetIsBanned() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IsBanned
}

func (p *Player) GetAvatar() avatar.Avatar {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Avatar
}

func (p *Player) GetAvatarChanged() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.AvatarChanged
}

func (p *Player) GetPos() simple.Position {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.position
}

func (p *Player) GetMastery() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Mastery.String()
}

func (p *Player) GetAchievements() simple.IntMap {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Achievements
}

func (p *Player) GetAchievementsLevel() simple.IntMap {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.achievementsLevel
}

func (p *Player) GetCoopID() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.CoopID
}

func (p *Player) GetIsStartedCoop() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IsStartedCoop
}

func (p *Player) GetIsSeekingJob() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.isSeekingJob
}

func (p *Player) GetLastLogin() time.Time {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.LastLogin
}

func (p *Player) GetDailyLogin() time.Time {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.DailyLogin
}

func (p *Player) GetGiftRefreshTime() time.Time {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.GiftRefreshTime
}

func (p *Player) GetGifts() gift.GiftList {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.Gifts
}

func (p *Player) GetIsRegistered() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IsRegistered
}

func (p *Player) GetIsTutorialCompleted() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.XP >= 10
}

func (p *Player) GetAccessLevel() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.AccessLevel
}

func (p *Player) GetMaxInstants() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.maxInstants
}

func (p *Player) GetJob() PlayerJob {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.job
}

func (p *Player) GetDishMasteryDuration(v int) int {
	masteryLevel := p.GetDishMasteryLevel(v)

	dishInfo, err := utils.GetDish(v)
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

func (p *Player) GetDishMasteryLevel(v int) int {
	// Get items info
	timesCooked := p.Mastery[v]
	dishInfo, err := utils.GetDish(v)
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

func (p *Player) GetDishMasteryServings(v int) int {
	masteryLevel := p.GetDishMasteryLevel(v)

	dishInfo, err := utils.GetDish(v)
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

func (p *Player) GetDishMasteryXP(v int) int {
	masteryLevel := p.GetDishMasteryLevel(v)

	dishInfo, err := utils.GetDish(v)
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

func (p *Player) GetIsInCoop() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.CoopID == 0 {
		return false
	}

	return true
}
