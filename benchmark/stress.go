package main

import (
	"bufio"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"math/rand"
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

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go mimicPlayer(&wg)
	}
	wg.Wait()
}

func createAccount(conn net.Conn, reader *bufio.Reader) string {
	name := RandString(10)
	email := name + "@bot.com"
	passwd := name

	//time.Sleep(10 * time.Millisecond)
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

func moveAround(conn net.Conn) {
	// Generate where to move
	x := rand.Intn(3) + 1
	xStr := strconv.Itoa(x)
	y := rand.Intn(3) + 1
	yStr := strconv.Itoa(y)

	// Move around ( "%xt%CafeEx%cwa%1%5%6%" )
	conn.Write([]byte("%xt%CafeEx%cwa%1%" + xStr + "%" + yStr + "%" + END_CHAR))
	time.Sleep(time.Duration((x-1)+(y-1)) * time.Second)
}

func mimicPlayer(wg *sync.WaitGroup) error {
	defer wg.Done()
	conn, err := net.Dial("tcp", "localhost:9339")
	if err != nil {
		println("Failed to create connection:", err.Error())
		return nil
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Register account
	id := createAccount(conn, reader)
	println(id)

	// Move around for 1 minute
	duration := 60 * time.Second
	endTime := time.Now().Add(duration)
	for time.Now().Before(endTime) {
		// Join Market
		conn.Write([]byte("%xt%CafeEx%mjm%1%" + END_CHAR))

		// Move around
		moveAround(conn)

		// Join own cafe ("%xt%CafeEx%jca%1%1%1%")
		// conn.Write([]byte("%xt%CafeEx%jca%1%" + id + "%" + id + "%" + END_CHAR))

		// Move around
		// moveAround(conn)
	}

	return nil
}
