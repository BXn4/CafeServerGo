package managers

import (
	"math/rand"
	"slices"
	"time"

	"github.com/charmbracelet/log"
)

func (gm *GameManager) CheckForShopAvailablity() {
	for true {
		currentTime := time.Now().UTC()
		start, end := GetShopUnavailabilityRange()
		duration := end.Sub(currentTime) //duration is in nanosec!
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

			gm.unavailableIngredients = choices
		}

		log.Infof("Shop will be unavailable from %s to %s", start, end)
		log.Infof("Unabailable ingredients will be: %d", gm.GetUnavailableIngredients())

		time.Sleep(duration) // duration is in nanosec!
		// sleep until end, so its resets the ingredients.
	}
}

func (gm *GameManager) GetUnavailableIngredients() []int {
	if len(gm.unavailableIngredients) != 0 {
		return gm.unavailableIngredients
	}

	return nil
}

func (gm *GameManager) IsIngredientUnavailable(ingredientID int) bool {
	if gm.IsShopUnavailable() && slices.Contains(gm.unavailableIngredients, ingredientID) {
		return true
	}

	return false
}

func (gm *GameManager) IsShopUnavailable() bool {
	currentTime := time.Now().UTC()
	start, end := GetShopUnavailabilityRange()
	return currentTime.After(start) && currentTime.Before(end)
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
