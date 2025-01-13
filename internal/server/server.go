package server

import (
	"cafego/internal/client"
	"cafego/internal/commands"
	"cafego/internal/database"
	"cafego/internal/managers"
	"cafego/internal/utils"
	"fmt"
	"log"
	"net"
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

	// Read the items XML file and cache it
	utils.ReadAndCacheItems()

	// Set up MariaDB connection
	db, err := database.ConnectToDB(s.dbConfig)
	if err != nil {
		println(err.Error())
		log.Fatal(err)
	}
	defer db.Close()
	log.Printf("Server connected to database...")

	s.gm.SetCafeDB(db)

	// Start the TCP server
	address := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Printf("Server started and listening on %s...", address)

	// Handle connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		c := client.New(conn, db)
		println("ADDING TO ClientManager")
		s.gm.AddClient(c)

		go commands.HandleClient(c, s.gm)
	}

}
