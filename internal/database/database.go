package database

import (
	"cafego/internal/objects"
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

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	// This only creates the db pbject and does not start a connection
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// We might want to use this in the future but for compatibility reasons we will not use it
	/*
		// err = db.AutoMigrate(&PlayerDAO{}, &CafeDAO{})
		// if err != nil {
		// 	return nil, err
		// }
	*/

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

func (db *CafeDB) CreateAccount(name, email, password string, avatar objects.Avatar) (*objects.Player, error) {

	hashedPasswd, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	dao := &PlayerDAO{
		Email:    email,
		Password: hashedPasswd,
		Username: name,
		Avatar:   avatar.String(name),
	}

	// Create player and get id
	err = db.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(dao).Error; err != nil {
			return fmt.Errorf("Cant create player: %w", err)
		}
		cafe := &CafeDAO{
			ID:        dao.ID,
			PlayerID:  dao.ID,
			OwnerName: name,
		}
		if err := tx.Create(cafe).Error; err != nil {
			return fmt.Errorf("cant create cafe: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Cant create cafe: %v", err)
	}

	// Parse player
	player, err := ConvertPlayerDAOToPlayer(*dao)
	if err != nil {
		return nil, fmt.Errorf("Cant parse player: %v", err)
	}

	return player, nil
}
