package database

import (
	"cafego/internal/models/player"
	"fmt"
	"strconv"
	"strings"

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

	// Build friends
	friendsStr := []string{}
	for _, f := range p.Friends {
		friendsStr = append(friendsStr, strconv.Itoa(f))
	}

	updateData := map[string]interface{}{
		"cash":                  int(p.GetCash()),
		"gold":                  p.GetGold(),
		"xp":                    p.GetXP(),
		"instant_cookings":      p.InstantCookings,
		"open_jobs":             p.OpenJobs,
		"coop_id":               p.CoopID,
		"played_wheel":          p.PlayedWheel,
		"allow_friend_requests": p.AllowFriendRequests,
		"friends":               strings.Join(friendsStr, "#"),
		"avatar":                p.Avatar.Apperance(),
		"mastery":               p.Mastery.String(),
		"achievement":           p.GetAchivements().String(),
		"gifts":                 p.Gifts.String(),
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
