/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package database

import (
	"cafego/internal/models/player"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (db *CafeDB) Authenticate(name string, pass string) (*player.Player, int, error) {
	// Find player by username or email
	var p player.Player
	err := db.conn.Where("username = ? OR email = ?", name, name).First(&p).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 14, fmt.Errorf("Player \"%v\" not found!", name)
		}
		return nil, 14, fmt.Errorf("DB Error: %v", err)
	}

	if !VerifyPassword(p.GetPassword(), pass) {
		return nil, 14, errors.New("Access Denied!")
	}

	if p.GetIsBanned() {
		return nil, 19, errors.New("You are banned!")
	}

	return &p, 0, nil
}

func (db *CafeDB) ChangePassword(id int, oldPass, newPass string) (int, error) {
	println("ChangePassword")
	// Find player by username or email
	var player player.Player
	err := db.conn.Where("id = ?", id).First(&player).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, fmt.Errorf("Player \"%v\" not found!", id)
		}
		return 1, fmt.Errorf("DB Error: %v", err)
	}

	if !VerifyPassword(player.GetPassword(), oldPass) {
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
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
