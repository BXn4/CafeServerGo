package utils

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Wod struct {
	ID               int    `xml:"id,attr"`
	Name             string `xml:"n,attr"`
	Group            string `xml:"g,attr"`
	Type             string `xml:"t,attr"`
	Walkable         int    `xml:"walkable,attr,omitempty"`
	Cash             int    `xml:"cash,attr,omitempty"`
	Gold             int    `xml:"gold,attr,omitempty"`
	Friends          int    `xml:"friends,attr,omitempty"`
	GoldNoFriends    int    `xml:"goldNoFriends,attr,omitempty"`
	Events           int    `xml:"events,attr,omitempty"`
	XP               int    `xml:"xp,attr,omitempty"`
	Level            int    `xml:"level,attr,omitempty"`
	InventorySize    int    `xml:"inventorySize,attr,omitempty"`
	IncomePerServing int    `xml:"incomePerServing,attr,omitempty"`
	Servings         int    `xml:"servings,attr,omitempty"`
	Duration         int    `xml:"duration,attr,omitempty"`
	DishCategory     int    `xml:"dishcategory,attr,omitempty"`
	Category         string `xml:"category,attr,omitempty"`
	Requirements     string `xml:"requirements,attr,omitempty"`
	ExpansionID      int    `xml:"expansionID,attr,omitempty"`
	SizeX            int    `xml:"sizeX,attr,omitempty"`
	SizeY            int    `xml:"sizeY,attr,omitempty"`
	MaxMembers       int    `xml:"maxMember,attr,omitempty"`
	MaxLevel         int    `xml:"maxLevel,attr,omitempty"`
	Chips            int    `xml:"chips,attr,omitempty"`
	RatingBonus      int    `xml:"ratingBonus,attr,omitempty"`
	Dishes           string `xml:"dishes,attr,omitempty"`
}

type Data struct {
	Wods []Wod `xml:"wod"`
}

// !!! THIS SHOUL ONLY BE READ !!!
var itemCollection map[string][]Wod

func ReadAndCacheItems() error {
	xmlFile, err := os.Open("./data/CafeItems.xml")
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer xmlFile.Close()
	fmt.Println("Successfully Opened CafeItems.xml")

	var result Data
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&result)
	if err != nil {
		return fmt.Errorf("Error decoding XML: %v\n", err)
	}

	itemCollection = make(map[string][]Wod)
	var loadedCount int
	for _, wod := range result.Wods {
		itemCollection[strings.ToLower(wod.Group)] = append(itemCollection[strings.ToLower(wod.Group)], wod)
		loadedCount++
	}

	fmt.Printf("Successfully loaded %d WOD entries\n", loadedCount)
	fmt.Println("Filtered and grouped WOD entries")
	return nil
}

func GetItems(s string) ([]Wod, error) {
	category := s
	if s == "fancy" {
		category = "ingredient"
	}
	items, ok := itemCollection[strings.ToLower(category)]
	if !ok {
		return nil, fmt.Errorf("Item category with name %v not exist", category)

	}

	// Filter items for fancy
	if s == "fancy" {
		newItems := []Wod{}
		for _, w := range items {
			if w.Category == "fancy" {
				newItems = append(newItems, w)
			}
		}
		items = newItems
	}

	return items, nil
}

func GetItem(id int) (Wod, error) {
	for _, wods := range itemCollection {
		for _, wod := range wods {
			if id == wod.ID {
				return wod, nil
			}
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetIngredient(id int) (Wod, error) {
	for _, item := range itemCollection["ingredient"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetFancyIngredient(id int) (Wod, error) {
	for _, item := range itemCollection["ingredient"] {
		if id == item.ID && item.Category == "fancy" {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetTile(id int) (Wod, error) {
	for _, item := range itemCollection["tile"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetWall(id int) (Wod, error) {
	for _, item := range itemCollection["wall"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDoor(id int) (Wod, error) {
	for _, item := range itemCollection["door"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetStove(id int) (Wod, error) {
	for _, item := range itemCollection["stove"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetCounter(id int) (Wod, error) {
	for _, item := range itemCollection["counter"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetFridge(id int) (Wod, error) {
	for _, item := range itemCollection["fridge"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetTable(id int) (Wod, error) {
	for _, item := range itemCollection["table"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDeco(id int) (Wod, error) {
	for _, item := range itemCollection["deco"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetChair(id int) (Wod, error) {
	for _, item := range itemCollection["chair"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetWallobject(id int) (Wod, error) {
	for _, item := range itemCollection["wallobject"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDish(id int) (Wod, error) {
	for _, item := range itemCollection["dish"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetVendingmachine(id int) (Wod, error) {
	for _, item := range itemCollection["vendingmachine"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetExpansion(id int) (Wod, error) {
	for _, item := range itemCollection["expansion"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetAchievement(id int) (Wod, error) {
	for _, item := range itemCollection["achievement"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetCoop(id int) (Wod, error) {
	for _, item := range itemCollection["coop"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetFastfood(id int) (Wod, error) {
	for _, item := range itemCollection["fastfood"] {
		if id == item.ID {
			return item, nil
		}
	}
	return Wod{}, fmt.Errorf("No item found with id: %v", id)
}
