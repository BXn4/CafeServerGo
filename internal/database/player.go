package database

import (
	"cafego/internal/models/player"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (db *CafeDB) GetPlayerByName(name string) (*player.Player, error) {

	var p player.Player
	err := db.conn.Where("username = ? OR email = ?", name, name).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("NAME: %v NOT FOUND", name)
		}
		return nil, fmt.Errorf("SQL ERR: %v", err)
	}

	p.Avatar.Name = p.Username

	return &p, nil

}

func (db *CafeDB) GetPlayer(id int) (*player.Player, error) {
	println("GetPlayer")
	var p player.Player
	err := db.conn.First(&p, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID NOT FOUND")
		}
		return nil, fmt.Errorf("SQL ERR: %v", err)
	}
	p.Avatar.Name = p.Username
	return &p, nil
}

func (db *CafeDB) SavePlayer(p *player.Player) error {

	p.LastLogin = time.Now().UTC()

	// Build friends
	friendsStr := []string{}
	for _, f := range p.Friends {
		friendsStr = append(friendsStr, strconv.Itoa(f))
	}

	friendsWithGiftsStr := []string{}
	for _, f := range p.FriendsWithGifts {
		friendsWithGiftsStr = append(friendsWithGiftsStr, strconv.Itoa(f))
	}

	updateData := map[string]any{
		"cash":                  int(p.GetCash()),
		"gold":                  p.GetGold(),
		"xp":                    p.GetXP(),
		"instant_cookings":      p.InstantCookings,
		"open_jobs":             p.OpenJobs,
		"coop_id":               p.CoopID,
		"played_wheel":          p.PlayedWheel,
		"allow_friend_requests": p.AllowFriendRequests,
		"friends":               strings.Join(friendsStr, "#"),
		"friends_with_gifts":    strings.Join(friendsWithGiftsStr, "#"),
		"sendable_gifts":        p.SendableGifts,
		"allow_emails":          p.AllowEmails,
		"email_verified":        p.EmailVerified,
		"username":              p.Username,
		"is_banned":             p.IsBanned,
		"avatar":                p.Avatar.Apperance(),
		"avatar_changed":        p.AvatarChanged,
		"mastery":               p.Mastery.String(),
		"achievement":           p.GetAchivements().String(),
		"last_login":            p.LastLogin,
		"daily_login":           p.DailyLogin,
		"gifts":                 p.Gifts.String(),
		"is_registered":         p.IsRegistered,
	}

	err := db.conn.Model(&player.Player{}).Where("id = ?", p.ID).Updates(updateData).Error
	if err != nil {
		return fmt.Errorf("Cant save player: %v", err)
	}
	return nil
}

func (db *CafeDB) DeleteFriend(playerID, friendID int) error {
	var p player.Player
	err := db.conn.First(&p, playerID).Error
	if err != nil {
		return fmt.Errorf("Player not found: %v", err)
	}

	index := -1
	for i, f := range p.Friends {
		if f == friendID {
			index = i
			break
		}
	}
	if index == -1 {
		return nil // Friend not found, no action needed
	}

	p.Friends = append(p.Friends[:index], p.Friends[index+1:]...)

	err = db.conn.Save(&p).Error
	if err != nil {
		return fmt.Errorf("Cant save friends: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateCash(playerID, playerCash int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("cash", playerCash).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateGold(playerID, playerGold int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("gold", playerGold).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) SetRegistered(playerID int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("is_registered", true).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateAvatarChanged(playerID int, playerAvatarChanged bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("avatar_changed", playerAvatarChanged).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateAvatar(playerID int, avatar string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("avatar", avatar).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateEmail(playerID int, email string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("email", email).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdatePassord(playerID int, password string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("password", password).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateXP(playerID int, xp int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("xp", xp).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateInstantCookings(playerID int, instantCookings int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("instant_cookings", instantCookings).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateOpenJobs(playerID int, openJobs int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("open_jobs", openJobs).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdatePlayedWheel(playerID int, playedWheel bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("played_wheel", playedWheel).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateAllowFriendRequests(playerID int, allowFriendRequests bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("allow_friend_requests", allowFriendRequests).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateFriends(playerID int, friends string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("friends", friends).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateFriendsWithGifts(playerID int, friendsWithGifts string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("friends_with_gifts", friendsWithGifts).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateSendableGifts(playerID int, sendableGifts string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("sendable_gifts", sendableGifts).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateAllowEmails(playerID int, allowEmails bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("allow_emails", allowEmails).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateEmailVerified(playerID int, emailVerified bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("email_verified", emailVerified).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateUsername(playerID int, username string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("username", username).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateIsBanned(playerID int, isBanned bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("is_banned", isBanned).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateMastery(playerID int, mastery string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("mastery", mastery).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateAchievement(playerID int, achievement string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("achievement", achievement).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateCoopID(playerID int, coopID int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("coop_id", coopID).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateLastLogin(playerID int, lastLogin time.Time) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("last_login", lastLogin).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateDailyLogin(playerID int, dailyLogin time.Time) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("daily_login", dailyLogin).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateGiftRefreshTime(playerID int, giftRefreshTime time.Time) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("gift_refresh_time", giftRefreshTime).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}

func (db *CafeDB) UpdateGifts(playerID int, gifts string) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("gifts", gifts).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}
