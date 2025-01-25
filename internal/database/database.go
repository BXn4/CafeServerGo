package database

import (
	"cafego/internal/objects"
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type CafeDB struct {
	conn *sql.DB
	mu   sync.Mutex
}

func ConnectToDB(config *DBConfig) (*CafeDB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	// This only creates the db pbject and does not start a connection
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// This checks if the connection is alive
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	cafe_db := &CafeDB{conn: db}

	return cafe_db, nil
}

func (db *CafeDB) Close() {
	db.conn.Close()
}

func (db *CafeDB) CreateAccount(name, email, password string, avatar objects.Avatar) (*objects.Player, error) {

	var id int
	// Create player and get id
	err := db.conn.QueryRow("INSERT INTO player ( email, password, username, avatar) VALUES (?,?,?,?) RETURNING id",
		email,
		password,
		name,
		avatar.String(),
	).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("Cant create player: %v", err)
	}

	// Create cafe
	_, err = db.conn.Exec("INSERT INTO cafe ( id, player_id, owner_name) VALUES (?,?,?)",
		id,
		id,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("Cant create cafe: %v", err)
	}

	// Parse player
	player, err := db.GetPlayer(id)
	if err != nil {
		return nil, err
	}

	return player, nil
}
