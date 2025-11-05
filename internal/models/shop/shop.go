package shop

import (
	"math/rand"
	"slices"
	"time"

	"github.com/charmbracelet/log"
)

var unavailableIngredients []int
var unavailable = false

func CheckForShopAvailablity(d time.Duration) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	Check()

	for range ticker.C {
		Check()
	}
}

func Check() {
	currentTime := time.Now().UTC()
	start, end := GetShopUnavailabilityRange()

	if currentTime.Before(start) || currentTime.After(end) {
		log.Infof("Shop is available now, no ingredients are unavailable.")

		setUnavailable(false)
		return
	}

	if !IsShopUnavailable() {

		var choices []int

		ingredients := map[int]bool{}

		for i := 1301; i <= 1343; i++ {
			ingredients[i] = true
		}

		alwaysAvailable := []int{
			1314, 1307, 1308, 1318,
			1327, 1319, 1316, 1311,
			1337, 1341, 1340,
		}

		for _, id := range alwaysAvailable {
			ingredients[id] = false
		}

		for i := 0; i < rand.Intn(4)+1; i++ {
			for {
				randomID := rand.Intn(1343-1301+1) + 1301
				if ingredients[randomID] {
					choices = append(choices, randomID)
					break
				}
			}

			setUnavailable(true)
			setUnavailableIngredietns(choices)
		}

		log.Infof("Shop is unavailable from %s to %s", start, end)
		log.Infof("Unabailable ingredients: %d", GetUnavailableIngredients())
	}
}

func GetUnavailableIngredients() []int {
	if len(unavailableIngredients) != 0 {
		return unavailableIngredients
	}

	return nil
}

func setUnavailableIngredietns(value []int) {
	unavailableIngredients = value
}

func setUnavailable(value bool) {
	unavailable = value
}

func IsIngredientUnavailable(ingredientID int) bool {
	if slices.Contains(unavailableIngredients, ingredientID) {
		return true
	}

	return false
}

func IsShopUnavailable() bool {
	return unavailable
}

func GetShopUnavailabilityRange() (time.Time, time.Time) {
	now := time.Now().UTC()
	seed := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	r := rand.New(rand.NewSource(seed))

	randomDayOffset := r.Intn(14)

	max := 12 // h
	min := 8  // h

	randomHourDuration := time.Duration(rand.Intn(max-min) + min)
	randomDate := now.AddDate(0, 0, randomDayOffset)

	startHour := r.Intn(22)
	start := time.Date(randomDate.Year(), randomDate.Month(), randomDate.Day(), startHour, 0, 0, 0, time.UTC)
	end := start.Add(randomHourDuration * time.Hour)

	return start, end
}
