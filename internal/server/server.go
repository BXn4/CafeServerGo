package server

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/database"
	"cafego/internal/managers"
	"cafego/internal/utils"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
)

// Config
type CafeConfig struct {
	Host string
	Port string
}

// Server
type CafeServer struct {
	config   *CafeConfig
	dbConfig *database.DBConfig
	maxConn  int
	db       *database.CafeDB
	gm       *managers.GameManager
}

func New(config *CafeConfig, dbconfig *database.DBConfig) (*CafeServer, error) {
	gm, err := managers.NewGameManager()
	if err != nil {
		return nil, err
	}
	return &CafeServer{
		config:   config,
		dbConfig: dbconfig,
		gm:       gm,
	}, nil
}

func (s *CafeServer) Run() {

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Read the items XML file and cache it
	utils.ReadAndCacheItems()

	// Read the levels XML file and cache it
	utils.ReadAndCacheLevels()

	// Read the achievements XML file and cache it
	utils.ReadAndCacheAchievements()

	// Set up MariaDB connection
	db, err := database.ConnectToDB(s.dbConfig)
	if err != nil {
		return
	}
	defer db.Close()
	log.Infof("Server connected to database.")

	s.gm.SetCafeDB(db)

	// Start the TCP server
	address := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Infof("Server started and listening on %s", address)

	// Handle connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			defer conn.Close()

			c := client.New(conn, db, s.gm)
			c.SetClientID(s.gm.NextClientID())
			s.gm.AddClient(c)
			c.Start()
			go commands.HandleClient(c, s.gm)
		}
	}()
	// Wait for interrupt signal
	<-sigChan
	log.Info("Received interrupt signal. Saving all data...")

	// Save all
	if err := s.gm.SaveAll(); err != nil {
		log.Fatalf("Error saving data: %v", err)
	} else {
		log.Info("All data saved successfully")
	}

}
