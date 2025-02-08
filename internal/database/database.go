package database

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/cafe"
	"cafego/internal/models/coop"
	"cafego/internal/models/player"
	"fmt"

	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type CafeDB struct {
	conn *gorm.DB
}

func ConnectToDB(config *DBConfig) (*CafeDB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true", config.User, config.Password, config.Host, config.Port, config.Database)

	// Creates the db pbject and does not start a connection
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Check if database tables are valid
	err = db.AutoMigrate(&player.Player{}, &cafe.Cafe{}, &coop.Coop{})
	if err != nil {
		return nil, err
	}

	cafe_db := &CafeDB{conn: db}

	return cafe_db, nil
}

func (db *CafeDB) Close() error {
	sqlDB, err := db.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (db *CafeDB) CreateAccount(name, email, password string, a avatar.Avatar) (*player.Player, error) {

	hashedPasswd, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	player := &player.Player{
		Email:    email,
		Password: hashedPasswd,
		Username: name,
		Avatar:   a,
	}

	// Create player and get id
	err = db.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(player).Error; err != nil {
			return fmt.Errorf("Cant create player: %w", err)
		}
		cafe := cafe.NewCafeForCreation(player.ID, player.ID, name)

		if err := tx.Create(cafe).Error; err != nil {
			return fmt.Errorf("Cant create cafe: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return player, nil
}
