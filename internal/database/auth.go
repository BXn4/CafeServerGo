package database

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func (db *CafeDB) Authenticate(name string, pass string) (int, error) {

	// Find player by username or email
	var player PlayerDAO
	err := db.conn.Where("username = ? OR email = ?", name, name).First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 14, fmt.Errorf("Player \"%v\" not found!", name)
		}
		return 14, fmt.Errorf("DB Error: %v", err)
	}

	// TODO: Secure authentication
	if player.Password != pass {
		return 14, errors.New("Access Denied!")
	}

	if player.IsBanned {
		return 19, errors.New("You are banned!")
	}

	return 0, nil
}
