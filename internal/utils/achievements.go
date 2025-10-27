package utils

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

type AchievementsWod struct {
	ID         int `xml:"a,attr"`
	WodID      int `xml:"w,attr"`
	Level      int `xml:"l,attr"`
	Target     int `xml:"r,attr"` // in game its called count, but target better for us.
	XPReward   int `xml:"x,attr"`
	CashReward int `xml:"ch,attr"`
	GoldReward int `xml:"g,attr"`
}

type AchievementRewards struct {
	Rewards []AchievementsWod `xml:"achievements"`
}

var achievementCollection map[int][]AchievementsWod

func ReadAndCacheAchievements() error {
	xmlFile, err := os.Open("./data/CafeAchievement.xml")
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer xmlFile.Close()
	log.Infof("Successfully Opened CafeAchievement.xml")

	var result AchievementRewards
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("Error decoding XML: %v", err)
	}

	achievementCollection = make(map[int][]AchievementsWod)
	var loadedCount int
	for _, rewardData := range result.Rewards {
		achievementCollection[rewardData.WodID] = append(achievementCollection[rewardData.WodID], rewardData)
		loadedCount++
	}

	log.Infof("Successfully loaded %d Achievement entries", loadedCount)
	return nil
}

func GetAchievementList(wodID int) []AchievementsWod {
	return achievementCollection[wodID]
}

func GetAchievementIDByWOD(wodID int) int {
	achievements, found := achievementCollection[wodID]
	if found {
		for _, achievement := range achievements {
			if achievement.WodID == wodID {
				return achievement.ID
			}
		}
	}

	return -1
}

func GetAchievementWODByID(id int) int {
	for wodID, list := range achievementCollection {
		for _, achievement := range list {
			if achievement.ID == id {
				return wodID
			}
		}
	}
	return -1
}

func GetAchievementData(id int, level int) (AchievementsWod, error) {
	rewardDatas, found := achievementCollection[id]
	if found {
		for _, achievement := range rewardDatas {
			if achievement.Level == level {
				return achievement, nil
			}
		}
	}

	return AchievementsWod{}, nil // Player is above the level, no more rewards!
}

func GetAchievementMAXLevel(wodID int) int {
	achievements, found := achievementCollection[wodID]
	if !found || len(achievements) == 0 {
		return -1
	}

	maxLevel := 0
	for _, achievement := range achievements {
		if achievement.Level > maxLevel {
			maxLevel = achievement.Level
		}
	}

	return maxLevel
}

func GetAchievementTarget(id int, level int) int {
	achievement, err := GetAchievementData(id, level)
	if err != nil {
		fmt.Printf("No achievement target data found to achievement: %d for level: %d", id, level)
		return -1
	}

	return achievement.Target
}

func GetAchievementXPReward(id int, level int) int {
	achievement, err := GetAchievementData(id, level)
	if err != nil {
		fmt.Printf("No achievement XP reward data found to level: %d", level)
		return 0
	}

	return achievement.XPReward
}

func GetAchievementCashReward(id int, level int) int {
	achievement, err := GetAchievementData(id, level)
	if err != nil {
		fmt.Printf("No achievement cash reward data found to level: %d", level)
		return 0
	}

	return achievement.CashReward
}

func GetAchievementGoldReward(id int, level int) int {
	achievement, err := GetAchievementData(id, level)
	if err != nil {
		fmt.Printf("No achievement gold reward data found to level: %d", level)
		return 0
	}

	return achievement.GoldReward
}
