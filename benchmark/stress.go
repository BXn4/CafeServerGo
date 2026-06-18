package main

import (
	"bufio"
	"database/sql"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const END_CHAR = "\x00"

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// Connect to existing database
func connectDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// Get account ID from database by game_id
func getAccountID(db *sql.DB, gameID string) (int64, error) {
	var accountID int64
	err := db.QueryRow("SELECT id FROM player WHERE id = ?", gameID).Scan(&accountID)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}

func createAccount(conn net.Conn, reader *bufio.Reader) string {
	name := RandString(10)
	email := name + "@bot.com"
	passwd := name

	conn.Write([]byte("<msg t='sys'><body action='verChk' r='0'><ver v='161' /></body></msg>" + END_CHAR))
	conn.Write([]byte("<msg t='sys'><body action='login' r='0'><login z='CafeEx'><nick><![CDATA[]]></nick><pword><![CDATA[201110190912%hu%null]]></pword></login></body></msg>" + END_CHAR))
	conn.Write([]byte("<msg t='sys'><body action='autoJoin' r='-1'></body></msg>" + END_CHAR))
	conn.Write([]byte("<msg t='sys'><body action='roundTrip' r='1'></body></msg>" + END_CHAR))
	conn.Write([]byte("%xt%CafeEx%vck%1%1603%" + END_CHAR))
	conn.Write([]byte("%xt%CafeEx%pin%1%" + END_CHAR))
	conn.Write([]byte("%xt%CafeEx%lca%1%" + name + "+2+1052$10#1062$0#1042$12#1082$0#1002$0#1022$1%1%63%1%<RoundHouseKick>%" + END_CHAR))
	conn.Write([]byte("%xt%CafeEx%lre%1%" + name + "%" + email + "%" + passwd + "%1%0%cafe%63%1%<RoundHouseKick>%-1%-1%-1%-1%-1%-1%-1%-1%" + END_CHAR))

	// Read until start with gui ( "%xt%gui%-1%0%6%6%" )
	var id string
	for {
		msg, err := reader.ReadString('\x00')
		if err != nil {
			return ""
		}
		if strings.HasPrefix(msg, "%xt%gui") {
			id = strings.Split(msg, "%")[5]
			break
		}
	}

	conn.Write([]byte("%xt%CafeEx%jca%1%" + id + "%" + id + "%" + END_CHAR))
	conn.Write([]byte("%xt%CafeEx%pin%1%" + END_CHAR))
	return id
}

func moveAround(conn net.Conn, db *sql.DB, accountID int64) {
	x := rand.Intn(3) + 1
	xStr := strconv.Itoa(x)
	y := rand.Intn(3) + 1
	yStr := strconv.Itoa(y)

	// Move around
	conn.Write([]byte("%xt%CafeEx%cwa%1%" + xStr + "%" + yStr + "%" + END_CHAR))

	// Log movement to database

	time.Sleep(time.Duration((x-1)+(y-1)) * time.Second)
}

func mimicPlayer(wg *sync.WaitGroup, db *sql.DB) error {
	defer wg.Done()

	conn, err := net.Dial("tcp", "localhost:9339")
	if err != nil {
		log.Printf("Failed to create connection: %v", err)
		return nil
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Register account (server handles this)
	id := createAccount(conn, reader)
	if id == "" {
		log.Printf("Failed to create account")
		return nil
	}
	log.Printf("Logged in with account ID: %s", id)

	// Get account ID from database if DB is available
	var accountID int64
	if db != nil {
		accountID, err = getAccountID(db, id)
		if err != nil {
			log.Printf("Account not found in DB (server may not have written yet): %v", err)
			// Continue without logging to DB
			accountID = 0
		} else {
		}
	}

	// Move around for 1 minute
	duration := 60 * time.Second
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		// Join Market
		conn.Write([]byte("%xt%CafeEx%mjm%1%" + END_CHAR))

		// Move around
		moveAround(conn, db, accountID)

		// Move around again
		moveAround(conn, db, accountID)
	}

	return nil
}

func main() {
	// Connect to existing database (optional)
	// If no database path is provided, the bot will run without logging
	dbPath := "../database/gg_cafe.db" // Change this to your actual database path

	var db *sql.DB
	var err error

	db, err = connectDB(dbPath)
	if err != nil {
		log.Printf("Warning: Could not connect to database at %s: %v", dbPath, err)
		log.Printf("Running without database logging...")
		db = nil
	} else {
		defer db.Close()
		log.Printf("Connected to database: %s", dbPath)
	}

	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		time.Sleep(1000 * time.Millisecond)
		wg.Add(1)
		go mimicPlayer(&wg, db)
	}
	wg.Wait()

	// Print statistics if database is available
	if db != nil {
	}
}
