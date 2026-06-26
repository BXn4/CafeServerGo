/*
CafeServerGo
A custom TCP socket server hosting library / game server.
Copyright (C) 2026 BXn4 and Hurka5
*/

package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type LevelWod struct {
	Level      int `xml:"l,attr"`
	Fridges    int `xml:"f,attr"`
	Counters   int `xml:"c,attr"`
	Instants   int `xml:"i,attr"`
	Stoves     int `xml:"s,attr"`
	Waiters    int `xml:"w,attr"`
	CashReward int `xml:"ch,attr"`
	GoldReward int `xml:"g,attr"`
}

type LevelLimitData struct {
	Levels []LevelWod `xml:"limit"`
}

var levelCollection map[int][]LevelWod

func ReadAndCacheLevels() error {
	xmlFile, err := os.Open("./data/CafeLevelXp.xml")
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer xmlFile.Close()
	log.Infof("Successfully Opened CafeLevelXp.xml")

	var result LevelLimitData
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("Error decoding XML: %v", err)
	}

	levelCollection = make(map[int][]LevelWod)
	var loadedCount int
	for _, levelData := range result.Levels {
		levelCollection[levelData.Level] = append(levelCollection[levelData.Level], levelData)
		loadedCount++
	}

	log.Infof("Successfully loaded %d Level entries", loadedCount)
	return nil
}

func GetLevelRewards(level int) (LevelWod, error) {
	for _, levelDatas := range levelCollection {
		for _, levelData := range levelDatas {
			if levelData.Level == level {
				return levelData, nil
			}
		}
	}

	return LevelWod{}, fmt.Errorf("No level data found to level: %v", level)
}

func GetLevelFridgesLimit(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No fridges data found to level: %v, using default: 1", level)
		return 1
	}
	return levelRewards.Fridges
}

func GetLevelCountersLimit(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No counters data found to level: %v, using default: 3", level)
		return 3
	}
	return levelRewards.Counters
}

func GetLevelInstantCookingsLimit(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No instant cooking data found to level: %v, using default: 12", level)
		return 12
	}
	return levelRewards.Instants
}

func GetLevelStovesLimit(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No stoves data found to level: %v, using default: 3", level)
		return 3
	}
	return levelRewards.Stoves
}

func GetLevelWaitersLimit(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No waiters data found to level: %v, using default: 1", level)
		return 1
	}
	return levelRewards.Waiters
}

func GetLevelCashReward(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No cash reward data found to level: %v, using default: 0", level)
		return 0
	}
	return levelRewards.CashReward
}

func GetLevelGoldReward(level int) int {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		fmt.Printf("No gold reward data found to level: %v, using default: 0", level)
		return 0
	}
	return levelRewards.GoldReward
}
