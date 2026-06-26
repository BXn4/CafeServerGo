/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package leaderboard

import (
	"cafego/internal/database"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type LeaderBoard struct {
	PlayerRank       int
	PlayerID         int
	PlayerXP         int
	PlayerCafeLuxury int
	PlayerName       string
}

const (
	OrderByXP int = iota
	OrderByLuxury
)

var (
	leaderboardXP     []LeaderBoard
	leaderboardLuxury []LeaderBoard
	mutex             sync.Mutex
)

func CheckLeaderBoard(db *database.CafeDB) {
	err := CacheLeaderBoard(db)
	if err != nil {
		log.Errorf("Failed to cache leaderboard! Err: %s", err)
	}

	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			err := CacheLeaderBoard(db)
			if err != nil {
				log.Errorf("Failed to cache leaderboard! Err: %s", err)
			}
		}
	}()
}

func CacheLeaderBoard(db *database.CafeDB) error {
	mutex.Lock()
	defer mutex.Unlock()

	list, err := db.GetLeaderBoard()
	if err != nil {
		log.Errorf("Failed to cache leaderboard! Err: %s", err)
	}

	leaderboardXP = make([]LeaderBoard, len(list))
	leaderboardLuxury = make([]LeaderBoard, len(list))

	for i, item := range list {
		leaderboard := LeaderBoard{
			PlayerRank:       item["rank"].(int),
			PlayerID:         item["id"].(int),
			PlayerXP:         item["xp"].(int),
			PlayerCafeLuxury: item["luxury"].(int),
			PlayerName:       item["username"].(string),
		}

		leaderboardXP[i] = leaderboard
		leaderboardLuxury[i] = leaderboard
	}

	sort.Slice(leaderboardXP, func(i, j int) bool {
		return leaderboardXP[i].PlayerXP > leaderboardXP[j].PlayerXP
	})

	for i := range leaderboardXP {
		leaderboardXP[i].PlayerRank = i + 1
	}

	sort.Slice(leaderboardXP, func(i, j int) bool {
		return leaderboardXP[i].PlayerRank < leaderboardXP[j].PlayerRank
	})

	sort.Slice(leaderboardLuxury, func(i, j int) bool {
		return leaderboardLuxury[i].PlayerCafeLuxury > leaderboardLuxury[j].PlayerCafeLuxury
	})

	for i := range leaderboardLuxury {
		leaderboardLuxury[i].PlayerRank = i + 1
	}

	sort.Slice(leaderboardLuxury, func(i, j int) bool {
		return leaderboardLuxury[i].PlayerRank < leaderboardLuxury[j].PlayerRank
	})

	log.Info("Leaderboard was updated!")
	return nil
}

// Some minor erros. not displaying the next page if the next page not contains 12 elemens
// Just for the scroll. the search works well.
func GetLeaderBoard(search string, orderBy int) string {
	mutex.Lock()
	defer mutex.Unlock()

	var leaderboard []LeaderBoard

	perPage := 12
	position := 0

	pos, err := strconv.Atoi(search)
	isSearchingPos := err == nil

	switch orderBy {
	case OrderByXP:
		leaderboard = leaderboardXP
	case OrderByLuxury:
		leaderboard = leaderboardLuxury
	default:
		leaderboard = leaderboardXP
	}

	for _, lb := range leaderboard {
		if (isSearchingPos && lb.PlayerRank == pos) ||
			(!isSearchingPos && lb.PlayerName == search) {
			position = lb.PlayerRank - 1
			break
		}
	}

	start := max(0, (position/perPage)*perPage)
	end := min(start+perPage, len(leaderboard))

	leaderboardList := leaderboard[start:end]

	var args []string

	for _, lb := range leaderboardList {
		leaderboardStr := fmt.Sprintf("%d+%d+%d+%s",
			lb.PlayerID, lb.PlayerXP, lb.PlayerCafeLuxury, lb.PlayerName)
		args = append(args, leaderboardStr)
	}

	return strings.Join(args, "#")
}
func GetPlayerRankByID(id, orderBy int) int {
	mutex.Lock()
	defer mutex.Unlock()

	var leaderboard []LeaderBoard

	switch orderBy {
	case OrderByXP:
		leaderboard = leaderboardXP
	case OrderByLuxury:
		leaderboard = leaderboardLuxury
	default:
		leaderboard = leaderboardXP
	}

	for _, lb := range leaderboard {
		if lb.PlayerID == id {
			return lb.PlayerRank - 1
		}
	}

	return 1
}

func GetPlayerRank(value string, orderBy int) int {
	mutex.Lock()
	defer mutex.Unlock()

	var leaderboard []LeaderBoard

	switch orderBy {
	case OrderByXP:
		leaderboard = leaderboardXP
	case OrderByLuxury:
		leaderboard = leaderboardLuxury
	default:
		leaderboard = leaderboardXP
	}

	position := 1

	pos, err := strconv.Atoi(value)
	isSearchingPos := err == nil

	for _, lb := range leaderboard {
		if (isSearchingPos && lb.PlayerRank == pos) ||
			(!isSearchingPos && lb.PlayerName == value) {
			position = lb.PlayerRank - 1
			break
		}
	}

	return position
}

// Some minor erros. not displaying the next page if the next page not contains 12 elemens
// Just for the scroll. the search works well.
func GetOffset(playerRank int) int {
	perPage := 12
	return max(0, (playerRank/perPage)*perPage)
}
