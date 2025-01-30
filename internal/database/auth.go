package database

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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

	if VerifyPassword(player.Password, pass) {
		return 14, errors.New("Access Denied!")
	}

	if player.IsBanned {
		return 19, errors.New("You are banned!")
	}

	return 0, nil
}

func (db *CafeDB) ChangePassword(id int, oldPass, newPass string) (int, error) {

	// Find player by username or email
	var player PlayerDAO
	err := db.conn.Where("id = ?", id).First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, fmt.Errorf("Player \"%v\" not found!", id)
		}
		return 1, fmt.Errorf("DB Error: %v", err)
	}

	if VerifyPassword(player.Password, oldPass) {
		return 10, errors.New("Access Denied!")
	}

	// Update password in database
	result := db.conn.Model(&player).Update("password", newPass)
	if result.Error != nil {
		return 1, fmt.Errorf("Failed to update password: %v", result.Error)
	}

	return 0, nil
}

// HashPassword generates a bcrypt hash for the given password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
