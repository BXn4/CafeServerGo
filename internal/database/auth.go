package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func (db *CafeDB) Authenticate(name string, pass string) (int, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	row := db.conn.QueryRow("SELECT password, is_banned, username FROM player WHERE username=? OR email=?", name, name)

	var password, username string
	var is_banned int
	err := row.Scan(
		&password,
		&is_banned,
		&username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return 14, fmt.Errorf("Player \"%v\" not found!", name)
		}
		return 14, fmt.Errorf("SQL ERR: %v", err)
	}

	// TODO: Secure authentication
	if password != pass {
		return 14, errors.New("Access Denied!")
	}

	if is_banned != 0 {
		return 19, errors.New("You are banned!")
	}

	return 0, nil
}
