package database

import (
	"cafego/internal/models/avatar"
	"cafego/internal/models/cafe"
	"cafego/internal/models/coops"
	"cafego/internal/models/player"
	"cafego/internal/models/simple"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBConfig struct {
	Database string
}

type CafeDB struct {
	conn *gorm.DB
}

func ConnectToDB(config *DBConfig) (*CafeDB, error) {
	db, err := gorm.Open(sqlite.Open("database/"+config.Database+".db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&player.Player{}, &cafe.Cafe{}, &coops.Coop{})
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

	achievements := make(simple.IntMap)
	for id := 2001; id <= 2030; id++ {
		achievements[id] = 0
	}

	mastery := make(simple.IntMap)
	for id := 1201; id <= 1255; id++ {
		mastery[id] = 0
	}

	player := &player.Player{}

	player.SetEmail(email)
	player.SetPassword(hashedPasswd)
	player.SetUsername(name)
	player.SetAvatar(a)
	player.SetAchievements(achievements)
	player.SetMastery(mastery)

	// Create player and get id
	err = db.conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(player).Error; err != nil {
			return fmt.Errorf("Cant create player: %w", err)
		}
		cafe := cafe.NewCafeForCreation(player.GetID(), player.GetID(), name)

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

func (db *CafeDB) GetLeaderBoard() ([]map[string]any, error) {
	type Entry struct {
		ID       int
		Username string
		XP       int
		Luxury   int
	}

	var rows []Entry
	err := db.conn.
		Table("player AS p").
		Select("p.id AS id, p.username AS username, p.xp AS xp, c.luxury AS luxury").
		Joins("LEFT JOIN cafe AS c ON c.owner_id = p.id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	leaderboard := make([]map[string]any, len(rows))
	for i, item := range rows {
		leaderboard[i] = map[string]any{
			"rank":     i + 1,
			"id":       item.ID,
			"username": item.Username,
			"xp":       item.XP,
			"luxury":   item.Luxury,
		}
	}

	return leaderboard, nil
}
