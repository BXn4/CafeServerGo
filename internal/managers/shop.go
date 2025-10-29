package managers

import (
	"cafego/internal/utils"
	"math/rand"
	"slices"
	"time"
)

func (gm *GameManager) CheckForShopAvailablity() {
	isShopUnavailable := true

	if isShopUnavailable && len(gm.unavailableIngredients) == 0 {
		ingredients, err := utils.GetItems("ingredients")
		if err != nil {
			return
		}
		var choices []int

		alwaysAvailable := []int{
			1314, 1307, 1308, 1318,
			1327, 1319, 1316, 1311,
			1337, 1341, 1340,
		}

		for i := 0; i < rand.Intn(4)+1; i++ {
			for {
				ingredient := ingredients[rand.Intn(len(ingredients))]
				if !slices.Contains(alwaysAvailable, ingredient.ID) {
					choices = append(choices, ingredient.ID)
					break
				}
			}
		}
		gm.unavailableIngredients = choices
	}
}

func (gm *GameManager) GetUnavailableIngredients() []int {
	if len(gm.unavailableIngredients) != 0 {
		return gm.unavailableIngredients
	}

	return nil
}

func (gm *GameManager) IsIngredientUnavailable(ingredientID int) bool {
	if slices.Contains(gm.unavailableIngredients, ingredientID) {
		return true
	}

	return false
}

func IsShopUnavailable(currentTime time.Time) bool {
	start, end := GetShopUnavailabilityRange()
	return currentTime.After(start) && currentTime.Before(end)
}

func GetShopUnavailabilityRange() (time.Time, time.Time) {
	now := time.Now().UTC()

	seed := int64(now.Year()*10000 + int(now.Month())*100 + now.Day())
	r := rand.New(rand.NewSource(seed))

	randomDayOffset := r.Intn(14)
	randomHourDuration := time.Duration(rand.Intn(4) + 1)
	randomDate := now.AddDate(0, 0, randomDayOffset)

	startHour := r.Intn(22)
	start := time.Date(randomDate.Year(), randomDate.Month(), randomDate.Day(), startHour, 0, 0, 0, time.UTC)
	end := start.Add(randomHourDuration * time.Hour)

	// fmt.Printf("Shop will be unavailable from %v to %v\n", start, end)
	return start, end
}
