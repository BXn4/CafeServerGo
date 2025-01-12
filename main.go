package main

import (
	"cafego/internal/database"
	"cafego/internal/server"
	"cafego/internal/utils"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	// Read .env file
	envFile, err := godotenv.Read(".env")

	hasConfig := err == nil

	if !hasConfig {
		fmt.Printf("Cannot find .env file!\n")
	}

	srv, err := server.New(
		// This is the server config
		&server.CafeConfig{
			Host: utils.If(hasConfig, envFile["SERVER_HOST"], "localhost"),
			Port: utils.If(hasConfig, envFile["SERVER_PORT"], "9339"),
		},
		// This is the database config
		&database.DBConfig{
			Host:     utils.If(hasConfig, envFile["DB_HOST"], "localhost"),
			Port:     utils.If(hasConfig, envFile["DB_PORT"], "3306"),
			Database: utils.If(hasConfig, envFile["DB_NAME"], "gg_cafe"),
			User:     utils.If(hasConfig, envFile["DB_USER"], "root"),
			Password: utils.If(hasConfig, envFile["DB_PASSWORD"], ""),
		},
	)

	if err != nil {
		fmt.Printf("Something went wrong with creating the server object!\n")
		os.Exit(1)
	} else {
		srv.Run()
	}
}
