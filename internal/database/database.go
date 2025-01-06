package database

import (
  _"log"
  "fmt"
  "database/sql"
	_ "github.com/go-sql-driver/mysql"
  "sync"
)

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password   string
    Database   string
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
  db.conn.Close();
}



