package main

import (
	"cafego/internal/database"
	"cafego/internal/models/balancing"
	"cafego/internal/server"
	"cafego/internal/utils"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

func init() {
	initStyles()

	// Uncomment to enable debug level logging:
	log.SetLevel(log.Level(-5))

	// Uncomment to enable info level logging:
	// log.SetLevel(log.InfoLevel)

}

func main() {

	// Read .env file
	envFile, err := godotenv.Read(".env")

	hasConfig := err == nil

	if !hasConfig {
		log.Warnf("Cannot find .env file!")
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

	balancing.LoadBalancing(hasConfig, envFile)

	if err != nil {
		log.Errorf("Failed to create the server object: %v", err)
		os.Exit(1)
	} else {
		srv.Run()
	}
}

func initStyles() {

	styles := log.DefaultStyles()

	// Default
	s := lipgloss.NewStyle().Bold(true)
	styles.Levels[log.ErrorLevel] = s.Copy().SetString("ERROR").Foreground(lipgloss.Color("#ff0000"))
	styles.Levels[log.WarnLevel] = s.Copy().SetString("WARN").Foreground(lipgloss.Color("#ffff00"))
	styles.Levels[log.InfoLevel] = s.Copy().SetString("INFO").Foreground(lipgloss.Color("#33ffcc"))
	styles.Levels[log.DebugLevel] = s.Copy().SetString("DEBUG").Foreground(lipgloss.Color("#7e9edf"))

	// Custom
	styles.Levels[log.Level(-3)] = s.Copy().SetString("SENT").Foreground(lipgloss.Color("#ffffed"))
	styles.Levels[log.Level(-2)] = s.Copy().SetString("ANNOUNCE").Foreground(lipgloss.Color("#fff3db"))
	styles.Levels[log.Level(-1)] = s.Copy().SetString("BROADCAST").Foreground(lipgloss.Color("#fff3db"))
	styles.Levels[log.Level(-5)] = s.Copy().SetString("RECEIVED").Foreground(lipgloss.Color("#c8eec8"))

	log.SetStyles(styles)
}
