package main

import (
	"cafego/internal/database"
	"cafego/internal/server"
)

func main() {
	srv := server.New(
		// This is the server config
		&server.CafeConfig{
			Host: "localhost",
			Port: "9339",
		},
		// This is the database config
		&database.DBConfig{
			Host:     "localhost",
			Port:     "3306",
			Database: "gg_cafe",
			User:     "root",
			Password: "1234",
		},
	)
	srv.Run()
}
