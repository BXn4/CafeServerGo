package database

import (
	"cafego/internal/objects"
	"database/sql"
	"fmt"
	_ "strconv"
	_ "strings"

	_ "github.com/go-sql-driver/mysql"
)

type PlayerDAO struct {
	ID                  int    `json:"id" form:"id" gorm:"column:id"`
	Email               string `json:"email" form:"email" gorm:"column:email"`
	Password            string `json:"password" form:"password" gorm:"column:password"`
	Cash                int    `json:"cash" form:"cash" gorm:"column:cash"`
	Gold                int    `json:"gold" form:"gold" gorm:"column:gold"`
	XP                  int    `json:"xp" form:"xp" gorm:"column:xp"`
	InstantCookings     int    `json:"instant_cookings" form:"instant_cookings" gorm:"column:instant_cookings"`
	OpenJobs            int    `json:"open_jobs" form:"open_jobs" gorm:"column:open_jobs"`
	PlayedWheel         bool   `json:"played_wheel" form:"played_wheel" gorm:"column:played_wheel"`
	AllowFriendRequests bool   `json:"allow_friend_requests" form:"allow_friend_requests" gorm:"column:allow_friend_requests"`
	AllowEmails         bool   `json:"allow_emails" form:"allow_emails" gorm:"column:allow_emails"`
	EmailVerified       bool   `json:"email_verified" form:"email_verified" gorm:"column:email_verified"`
	NewGifts            int    `json:"new_gifts" form:"new_gifts" gorm:"column:new_gifts"`
	Username            string `json:"username" form:"username" gorm:"column:username"`
	Gender              int    `json:"gender" form:"gender" gorm:"column:gender"`
	TopColor            int    `json:"top_color" form:"top_color" gorm:"column:top_color"`
	SkinColor           int    `json:"skin_color" form:"skin_color" gorm:"column:skin_color"`
	HairColor           int    `json:"hair_color" form:"hair_color" gorm:"column:hair_color"`
	LegsColor           int    `json:"legs_color" form:"legs_color" gorm:"column:legs_color"`
	IsBanned            bool   `json:"is_banned" form:"is_banned" gorm:"column:is_banned"`
	Mastery             string `json:"mastery" form:"mastery" gorm:"column:mastery"`
	LastLogin           string `json:"last_login" form:"last_login" gorm:"column:last_login"`
	DailyLogin          string `json:"daily_login" form:"daily_login" gorm:"column:daily_login"`
}

func ConvertPlayerDAOToPlayer(playerDAO PlayerDAO) (*objects.Player, error) {
	var player objects.Player

	// Fill in simple data
	player.ID = playerDAO.ID
	player.Cash = playerDAO.Cash
	player.Gold = playerDAO.Gold
	player.XP = playerDAO.XP
	player.InstantCookings = playerDAO.InstantCookings
	player.OpenJobs = playerDAO.OpenJobs
	player.PlayedWheel = playerDAO.PlayedWheel
	player.AllowFriendRequests = playerDAO.AllowFriendRequests
	player.AllowEmails = playerDAO.AllowEmails
	player.EmailVerified = playerDAO.EmailVerified
	player.NewGifts = playerDAO.NewGifts
	player.Username = playerDAO.Username
	player.Position = []int{0, 0}

	player.ParseMastery(playerDAO.Mastery)

	// Fill avatar
	avatar := objects.Avatar{
		Name:      playerDAO.Username,
		Gender:    objects.AvatarGender(playerDAO.Gender),
		SkinColor: playerDAO.SkinColor,
		TopColor:  playerDAO.TopColor,
		HairColor: playerDAO.HairColor,
		LegsColor: playerDAO.LegsColor,
		IsNPC:     false,
	}

	player.Avatar = avatar

	return &player, nil
}

func (db *CafeDB) GetPlayerByName(name string) (*objects.Player, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.conn.QueryRow("SELECT * FROM player WHERE username = ? OR email = ?", name, name)

	var playerDAO PlayerDAO
	err := row.Scan(
		&playerDAO.ID,
		&playerDAO.Email,
		&playerDAO.Password,
		&playerDAO.Cash,
		&playerDAO.Gold,
		&playerDAO.XP,
		&playerDAO.InstantCookings,
		&playerDAO.OpenJobs,
		&playerDAO.PlayedWheel,
		&playerDAO.AllowFriendRequests,
		&playerDAO.AllowEmails,
		&playerDAO.EmailVerified,
		&playerDAO.NewGifts,
		&playerDAO.Username,
		&playerDAO.Gender,
		&playerDAO.TopColor,
		&playerDAO.SkinColor,
		&playerDAO.HairColor,
		&playerDAO.LegsColor,
		&playerDAO.IsBanned,
		&playerDAO.Mastery,
		&playerDAO.LastLogin,
		&playerDAO.DailyLogin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
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
