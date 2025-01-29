package database

import (
	"cafego/internal/objects"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type PlayerDAO struct {
	ID                  int    `gorm:"column:id"`
	Email               string `gorm:"column:email"`
	Password            string `gorm:"column:password"`
	Cash                int    `gorm:"column:cash"`
	Gold                int    `gorm:"column:gold"`
	XP                  int    `gorm:"column:xp"`
	InstantCookings     int    `gorm:"column:instant_cookings"`
	OpenJobs            int    `gorm:"column:open_jobs"`
	PlayedWheel         bool   `gorm:"column:played_wheel"`
	AllowFriendRequests bool   `gorm:"column:allow_friend_requests"`
	Friends             string `gorm:"column:friends"`
	FriendsWithGifts    string `gorm:"column:friends_with_gifts"`
	AllowEmails         bool   `gorm:"column:allow_emails"`
	EmailVerified       bool   `gorm:"column:email_verified"`
	Username            string `gorm:"column:username"`
	Avatar              string `gorm:"column:avatar"`
	IsBanned            bool   `gorm:"column:is_banned"`
	Mastery             string `gorm:"column:mastery"`
	Achievement         string `gorm:"column:achievement"`
	LastLogin           string `gorm:"column:last_login"`
	DailyLogin          string `gorm:"column:daily_login"`
	Gifts               string `gorm:"column:gifts"`
	GiftRefreshTime     string `gorm:"column:gift_refresh_time"`
	SendableGifts       string `gorm:"column:sendable_gifts"`
	AccessLevel         int    `gorm:"column:access_level"`
}

func (playerDAO PlayerDAO) TableName() string {
	return "player"
}

func ConvertPlayerDAOToPlayer(playerDAO PlayerDAO) (*objects.Player, error) {

	var player objects.Player

	// Fill in simple data
	player.ID = playerDAO.ID
	player.SetCash(playerDAO.Cash)
	player.SetGold(playerDAO.Gold)
	player.SetXP(playerDAO.XP)
	player.ParseAchievement(playerDAO.Achievement)
	player.InstantCookings = playerDAO.InstantCookings
	player.OpenJobs = playerDAO.OpenJobs
	player.PlayedWheel = playerDAO.PlayedWheel
	player.AllowFriendRequests = playerDAO.AllowFriendRequests
	player.AllowEmails = playerDAO.AllowEmails
	player.EmailVerified = playerDAO.EmailVerified
	player.Username = playerDAO.Username
	player.AccessLevel = playerDAO.AccessLevel
	player.Position = [2]int{0, 0}

	// Parse gifts
	gifts, err := objects.ParseGifts(playerDAO.Gifts)
	if err != nil {
		return nil, err
	}
	player.Gifts = gifts

	player.ParseFriends(playerDAO.Friends)
	player.ParseMastery(playerDAO.Mastery)

	// Fill avatar
	player.Avatar = *objects.NewAvatarFromString(playerDAO.Avatar)
	player.Avatar.IsNPC = false

	return &player, nil
}

func (db *CafeDB) GetPlayerByName(name string) (*objects.Player, error) {

	var playerDAO PlayerDAO
	err := db.conn.Where("username = ? OR email = ?", name, name).First(&playerDAO).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("NAME: %v NOT FOUND", name)
		}
		return nil, fmt.Errorf("SQL ERR: %v", err)
	}

	player, err := ConvertPlayerDAOToPlayer(playerDAO)
	if err != nil {
		return nil, err
	}

	return player, nil

}

func (db *CafeDB) GetPlayer(id int) (*objects.Player, error) {
	var playerDAO PlayerDAO
	err := db.conn.First(&playerDAO, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID NOT FOUND")
		}
		return nil, fmt.Errorf("SQL ERR: %v", err)
	}

	player, err := ConvertPlayerDAOToPlayer(playerDAO)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func (db *CafeDB) SavePlayer(player *objects.Player) error {

	// Build friends
	friendsStr := []string{}
	for _, f := range player.Friends {
		friendsStr = append(friendsStr, strconv.Itoa(f))
	}

	updateData := map[string]interface{}{
		"cash":                  uint(player.GetCash()),
		"gold":                  player.GetGold(),
		"xp":                    player.GetXP(),
		"instant_cookings":      player.InstantCookings,
		"open_jobs":             player.OpenJobs,
		"played_wheel":          player.PlayedWheel,
		"allow_friend_requests": player.AllowFriendRequests,
		"friends":               strings.Join(friendsStr, "#"),
		"avatar":                player.Avatar.Apperance(),
		"mastery":               player.BuildMastery(),
		"achievement":           player.BuildAchievement(),
		"gifts":                 objects.BuildGifts(player.Gifts),
	}

	err := db.conn.Model(&PlayerDAO{}).Where("id = ?", player.ID).Updates(updateData).Error
	if err != nil {
		return fmt.Errorf("Cant save player: %v", err)
	}
	return nil
}

func (db *CafeDB) DeleteFriend(playerID, friendID int) error {
	var playerDAO PlayerDAO
	err := db.conn.First(&playerDAO, playerID).Error
	if err != nil {
		return fmt.Errorf("Player not found: %v", err)
	}

	friends := strings.Split(playerDAO.Friends, "#")
	friendIDStr := strconv.Itoa(friendID)
	index := -1
	for i, f := range friends {
		if f == friendIDStr {
			index = i
			break
		}
	}
	if index == -1 {
		return nil // Friend not found, no action needed
	}

	newFriends := append(friends[:index], friends[index+1:]...)
	playerDAO.Friends = strings.Join(newFriends, "#")

	err = db.conn.Save(&playerDAO).Error
	if err != nil {
		return fmt.Errorf("Cant save friends: %v", err)
	}
	return nil
}
func (db *CafeDB) GetDailyLogin(playerID int) (*time.Time, error) {
	var playerDAO PlayerDAO
	err := db.conn.Select("daily_login").First(&playerDAO, playerID).Error
	if err != nil {
		return nil, err
	}

	dailyLogin, err := time.Parse("2006-01-02 15:04:05", playerDAO.DailyLogin)
	if err != nil {
		return nil, fmt.Errorf("Error parsing daily login time: %v", err)
	}

	return &dailyLogin, nil
}

func (db *CafeDB) GetGiftRefreshTime(playerID int) (*time.Time, error) {
	var playerDAO PlayerDAO
	err := db.conn.Select("gift_refresh_time").First(&playerDAO, playerID).Error
	if err != nil {
		return nil, err
	}

	refreshTime, err := time.Parse("2006-01-02 15:04:05", playerDAO.GiftRefreshTime)
	if err != nil {
		return nil, fmt.Errorf("Error parsing daily login time: %v", err)
	}

	return &refreshTime, nil
}

func (db *CafeDB) GetFriendsWithGifts(playerID int) ([]string, error) {
	var playerDAO PlayerDAO
	err := db.conn.Select("friends_with_gifts").First(&playerDAO, playerID).Error
	if err != nil {
		return nil, err
	}

	return strings.Split(playerDAO.FriendsWithGifts, "+"), nil
}

func (db *CafeDB) GetSendableGifts(playerID int) (string, error) {
	var playerDAO PlayerDAO
	err := db.conn.Select("gift_refresh_time").First(&playerDAO, playerID).Error
	if err != nil {
		return "", err
	}

	return playerDAO.SendableGifts, nil
}

func (db *CafeDB) ResetDailyLogin(playerID int) error {
	err := db.conn.Model(&PlayerDAO{}).Where("id = ?", playerID).Update("daily_login", time.Now().Format("2006-01-02 15:04:05")).Error
	if err != nil {
		return fmt.Errorf("Cant reset daily login: %v", err)
	}
	return nil
}

func (db *CafeDB) ResetGiftRefreshTime(playerID int) error {
	err := db.conn.Model(&PlayerDAO{}).Where("id = ?", playerID).Update("gift_refresh_time", time.Now().Format("2006-01-02 15:04:05")).Error
	if err != nil {
		return fmt.Errorf("Cant reset gift refresh time: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateSendableGifts(playerID int, gifts string) error {
	err := db.conn.Model(&PlayerDAO{}).Where("id = ?", playerID).Update("sendable_gifts", gifts).Error
	if err != nil {
		return fmt.Errorf("Cant update sendable gifts: %v", err)
	}
	return nil
}
