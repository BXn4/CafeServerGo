package database

import (
	"cafego/internal/models/player"
	"fmt"
	"strconv"
	"strings"
	"time"

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

	return &p, nil

}

func (db *CafeDB) GetPlayer(id int) (*player.Player, error) {
	var p player.Player
	err := db.conn.First(&p, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID NOT FOUND")
		}
		return nil, fmt.Errorf("SQL ERR: %v", err)
	}

	return &p, nil
}

func (db *CafeDB) SavePlayer(p *player.Player) error {

	p.SetLastLogin(time.Now().UTC())

	// Build friends
	friendsStr := []string{}
	for _, f := range p.GetFriends() {
		friendsStr = append(friendsStr, strconv.Itoa(f))
	}

	friendsWithGiftsStr := []string{}
	for _, f := range p.GetFriendsWithGifts() {
		friendsWithGiftsStr = append(friendsWithGiftsStr, strconv.Itoa(f))
	}

	avatar := p.GetAvatar()

	updateData := map[string]any{
		"cash":                  int(p.GetCash()),
		"gold":                  p.GetGold(),
		"xp":                    p.GetXP(),
		"instant_cookings_used": p.GetInstantCookingsUsed(),
		"open_jobs":             p.GetOpenJobs(),
		"refilled_jobs":         p.GetRefilledJobs(),
		"coop_id":               p.GetCoopID(),
		"played_wheel":          p.GetPlayedWheel(),
		"allow_friend_requests": p.GetAllowFriendRequests(),
		"friends":               strings.Join(friendsStr, "#"),
		"friends_with_gifts":    strings.Join(friendsWithGiftsStr, "#"),
		"sendable_gifts":        p.GetSendableGifts(),
		"allow_emails":          p.GetAllowEmails(),
		"email_verified":        p.GetEmailVerified(),
		"username":              p.GetUsername(),
		"is_banned":             p.GetIsBanned(),
		"avatar":                avatar.Apperance(),
		"avatar_changed":        p.GetAvatarChanged(),
		"mastery":               p.GetMastery(),
		"achievements":          p.GetAchivements().String(),
		"last_login":            p.GetLastLogin(),
		"daily_login":           p.GetDailyLogin(),
		"gifts":                 p.GetGifts().String(),
	}

	err := db.conn.Model(&player.Player{}).Where("id = ?", p.GetID()).Updates(updateData).Error
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
	for i, f := range p.GetFriends() {
		if f == friendID {
			index = i
			break
		}
	}
	if index == -1 {
		return nil // Friend not found, no action needed
	}

	p.SetFriends(append(p.GetFriends()[:index], p.GetFriends()[index+1:]...))

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
	println(xp)
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

func (db *CafeDB) UpdateRefilledJobs(playerID int, refilledJobs int) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("refilled_jobs", refilledJobs).Error
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
		Update("achievements", achievement).Error
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

func (db *CafeDB) UpdateStartedCoop(playerID int, started_coop bool) error {
	err := db.conn.Model(&player.Player{}).
		Where("id = ?", playerID).
		Update("started_coop", started_coop).Error
	if err != nil {
		return fmt.Errorf("Cant update player: %v", err)
	}
	return nil
}
