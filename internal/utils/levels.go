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

func GetLevelFridgesLimit(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 1, fmt.Errorf("No fridges data found to level: %v, using default: 1", level)
	}
	return levelRewards.Fridges, nil
}

func GetLevelCountersLimit(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 3, fmt.Errorf("No counters data found to level: %v, using default: 3", level)
	}
	return levelRewards.Counters, nil
}

func GetLevelInstantCookingsLimit(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 12, fmt.Errorf("No instant cooking data found to level: %v, using default: 12", level)
	}
	return levelRewards.Instants, nil
}

func GetLevelStovesLimit(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 3, fmt.Errorf("No stoves data found to level: %v, using default: 3", level)
	}
	return levelRewards.Stoves, nil
}

func GetLevelWaitersLimit(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 1, fmt.Errorf("No waiters data found to level: %v, using default: 1", level)
	}
	return levelRewards.Waiters, nil
}

func GetLevelCashReward(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 0, fmt.Errorf("No cash reward data found to level: %v, using default: 0", level)
	}
	return levelRewards.CashReward, nil
}

func GetLevelGoldReward(level int) (int, error) {
	levelRewards, err := GetLevelRewards(level)
	if err != nil {
		return 0, fmt.Errorf("No gold reward data found to level: %v, using default: 0", level)
	}
	return levelRewards.GoldReward, nil
}
