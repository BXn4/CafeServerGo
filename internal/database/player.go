package database

import (
	"cafego/internal/objects"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

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
	Friends             string `json:"friends" form:"friends" gorm:"column:friends"`
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
	Achievement         string `json:"achievement" form:"achievement" gorm:"column:achievement"`
	LastLogin           string `json:"last_login" form:"last_login" gorm:"column:last_login"`
	DailyLogin          string `json:"daily_login" form:"daily_login" gorm:"column:daily_login"`
	Gifts               string `json:"gifts" form:"gifts" gorm:"column:gifts"`
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
	player.Position = [2]int{0, 0}

	// Parse gifts
	gifts, err := objects.ParseGifts(playerDAO.Gifts)
	if err != nil {
		return nil, err
	}
	player.Gifts = gifts

	player.ParseFriends(playerDAO.Friends)
	player.ParseMastery(playerDAO.Mastery)
	player.ParseAchievement(playerDAO.Achievement)

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

	// Prepare query
	row := db.conn.QueryRow("SELECT * FROM player WHERE username = ? OR email = ?", name, name)

	// Scan rows
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
		&playerDAO.Friends,
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
		&playerDAO.Achievement,
		&playerDAO.LastLogin,
		&playerDAO.DailyLogin,
		&playerDAO.Gifts,
	)
	if err != nil {
		if err == sql.ErrNoRows {
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
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.conn.QueryRow("SELECT * FROM player WHERE id = ?", id)

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
		&playerDAO.Friends,
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
		&playerDAO.Achievement,
		&playerDAO.LastLogin,
		&playerDAO.DailyLogin,
		&playerDAO.Gifts,
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

func (db *CafeDB) SavePlayer(player *objects.Player) {
	db.mu.Lock()
	defer db.mu.Unlock()

	friendsStr := []string{}
	for _, f := range player.Friends {
		friendsStr = append(friendsStr, strconv.Itoa(f))
	}

	result, err := db.conn.Exec(
		" UPDATE player SET "+
			"cash = ?,"+
			"gold = ?,"+
			"xp = ?,"+
			"instant_cookings = ?,"+
			"open_jobs = ?,"+
			"played_wheel = ?,"+
			"open_jobs = ?,"+
			"played_wheel = ?,"+
			"allow_friend_requests = ?,"+
			"friends = ?,"+
			"new_gifts = ?,"+
			"gender = ?,"+
			"top_color = ?,"+
			"skin_color = ?,"+
			"hair_color = ?,"+
			"legs_color = ?,"+
			"mastery = ?,"+
			"achievement = ?,"+
			"gifts = ? "+
			"WHERE id = ?",
		uint(player.Cash),
		player.Gold,
		player.XP,
		player.InstantCookings,
		player.OpenJobs,
		player.PlayedWheel,
		player.OpenJobs,
		player.PlayedWheel,
		player.AllowFriendRequests,
		strings.Join(friendsStr, "#"),
		player.NewGifts,
		player.Avatar.Gender,
		player.Avatar.TopColor,
		player.Avatar.SkinColor,
		player.Avatar.HairColor,
		player.Avatar.LegsColor,
		player.BuildMastery(),
		player.BuildAchievement(),
		objects.BuildGifts(player.Gifts),
		player.ID,
	)

	if err != nil {
		fmt.Printf("Cant save player: %v\n", err)
		return
	}

	// Check how many rows were affected
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("Error fetching rows affected:", err)
	}
}

func (db *CafeDB) DeleteFriend(playerID, friendID int) {
	db.mu.Lock()
	defer db.mu.Unlock()

	friendIDStr := strconv.Itoa(friendID)

	// Get friends
	row := db.conn.QueryRow("SELECT friends FROM player WHERE id = ?", playerID)
	var friendsRaw string
	err := row.Scan(&friendsRaw)

	// Delete friend
	friends := strings.Split(friendsRaw, "#")
	index := -1
	for i, f := range friends {
		if f == friendIDStr {
			index = i
		}
	}
	if index == -1 {
		return
	}

	newFriends := append(friends[:index], friends[index+1:]...)

	// Update friends
	result, err := db.conn.Exec(
		" UPDATE player SET friends = ? WHERE id = ?",
		strings.Join(newFriends, "#"),
		playerID,
	)

	if err != nil {
		fmt.Printf("Cant save friends: %v\n", err)
		return
	}

	// Check how many rows were affected
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("Error fetching rows affected:", err)
	}
}

func (db *CafeDB) GetDailyLogin(playerID int) (*time.Time, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.conn.QueryRow("SELECT daily_login FROM player WHERE id = ?", playerID)
	var dailyLoginStr string

	err := row.Scan(&dailyLoginStr)

	if err != nil {
		return nil, err
	}

	dailyLogin, err := time.Parse("2025-01-23 22:09:07", dailyLoginStr)

	return &dailyLogin, nil
}

func (db *CafeDB) ResetDailyLogin(playerID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Update friends
	result, err := db.conn.Exec(
		"UPDATE player SET daily_login = ? WHERE id = ?",
		time.Now(),
		playerID,
	)

	if err != nil {
		return fmt.Errorf("Cant reset daily_login: %v\n", err)
	}

	// Check how many rows were affected
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("Error fetching rows affected:", err)
	}

	return nil
}
