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
	id                    int    `json:"id" form:"id" gorm:"column:id"`
	email                 string `json:"email" form:"email" gorm:"column:email"`
	password              string `json:"password" form:"password" gorm:"column:password"`
	cash                  int    `json:"cash" form:"cash" gorm:"column:cash"`
	gold                  int    `json:"gold" form:"gold" gorm:"column:gold"`
	xp                    int    `json:"xp" form:"xp" gorm:"column:xp"`
	instant_cookings      int    `json:"instant_cookings" form:"instant_cookings" gorm:"column:instant_cookings"`
	open_jobs             int    `json:"open_jobs" form:"open_jobs" gorm:"column:open_jobs"`
	played_wheel          bool   `json:"played_wheel" form:"played_wheel" gorm:"column:played_wheel"`
	allow_friend_requests bool   `json:"allow_friend_requests" form:"allow_friend_requests" gorm:"column:allow_friend_requests"`
	allow_emails          bool   `json:"allow_emails" form:"allow_emails" gorm:"column:allow_emails"`
	email_verified        bool   `json:"email_verified" form:"email_verified" gorm:"column:email_verified"`
	new_gifts             int    `json:"new_gifts" form:"new_gifts" gorm:"column:new_gifts"`
	username              string `json:"username" form:"username" gorm:"column:username"`
	gender                int    `json:"gender" form:"gender" gorm:"column:gender"`
	top_color             int    `json:"top_color" form:"top_color" gorm:"column:top_color"`
	skin_color            int    `json:"skin_color" form:"skin_color" gorm:"column:skin_color"`
	hair_color            int    `json:"hair_color" form:"hair_color" gorm:"column:hair_color"`
	legs_color            int    `json:"legs_color" form:"legs_color" gorm:"column:legs_color"`
	is_banned             bool   `json:"is_banned" form:"is_banned" gorm:"column:is_banned"`
	mastery               string `json:"mastery" form:"mastery" gorm:"column:mastery"`
	last_login            string `json:"last_login" form:"last_login" gorm:"column:last_login"`
	daily_login           string `json:"daily_login" form:"daily_login" gorm:"column:daily_login"`
}

func ConvertPlayerDAOToPlayer(playerDAO PlayerDAO) (*objects.Player, error) {
	var player objects.Player

	// Fill in simple data
	player.ID = playerDAO.id
	player.Cash = playerDAO.cash
	player.Gold = playerDAO.gold
	player.XP = playerDAO.xp
	player.InstantCookings = playerDAO.instant_cookings
	player.OpenJobs = playerDAO.open_jobs
	player.PlayedWheel = playerDAO.played_wheel
	player.AllowFriendRequests = playerDAO.allow_friend_requests
	player.AllowEmails = playerDAO.allow_emails
	player.EmailVerified = playerDAO.email_verified
	player.NewGifts = playerDAO.new_gifts
	player.Username = playerDAO.username
	player.Position = []int{0, 0}

	player.ParseMastery(playerDAO.mastery)

	// Fill avatar
	avatar := objects.Avatar{
		playerDAO.username,
		objects.AvatarGender(playerDAO.gender),
		playerDAO.skin_color,
		playerDAO.top_color,
		playerDAO.hair_color,
		playerDAO.legs_color,
		false,
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
		&playerDAO.id,
		&playerDAO.email,
		&playerDAO.password,
		&playerDAO.cash,
		&playerDAO.gold,
		&playerDAO.xp,
		&playerDAO.instant_cookings,
		&playerDAO.open_jobs,
		&playerDAO.played_wheel,
		&playerDAO.allow_friend_requests,
		&playerDAO.allow_emails,
		&playerDAO.email_verified,
		&playerDAO.new_gifts,
		&playerDAO.username,
		&playerDAO.gender,
		&playerDAO.top_color,
		&playerDAO.skin_color,
		&playerDAO.hair_color,
		&playerDAO.legs_color,
		&playerDAO.is_banned,
		&playerDAO.mastery,
		&playerDAO.last_login,
		&playerDAO.daily_login,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Errorf("ID NOT FOUND")
		}
		fmt.Errorf("SQL ERR: %v", err)
		return nil, err
	}

	player, err := ConvertPlayerDAOToPlayer(playerDAO)
	if err != nil {
		return nil, err
	}

	return player, nil
}
