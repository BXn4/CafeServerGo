package utils

import (
	"encoding/xml"
	"fmt"
	"os"
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
	Duration         int    `xml:"Duration,attr,omitempty"`
	Category         int    `xml:"dishcategory,attr,omitempty"`
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
		fmt.Println("Error opening file:", err)
		return err
	}
	defer xmlFile.Close()
	fmt.Println("Successfully Opened CafeItems.xml")

	var result Data
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return err
	}

  itemCollection = make(map[string][]Wod)
  var loadedCount int
	for _, wod := range result.Wods {
    itemCollection[wod.Group] = append(itemCollection[wod.Group], wod)
    loadedCount++
  }

	fmt.Printf("Successfully loaded %d WOD entries\n", loadedCount)
	fmt.Println("Filtered and grouped WOD entries")
  return nil
}


func GetIngredient(id int) (Wod, error) {
	for _, item := range itemCollection["Ingredient"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}


func GetTile(id int) (Wod, error) {
	for _, item := range itemCollection["Tile"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetWall(id int) (Wod, error) {
	for _, item := range itemCollection["Wall"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDoor(id int) (Wod, error) {
	for _, item := range itemCollection["Door"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetStove(id int) (Wod, error) {
	for _, item := range itemCollection["Stove"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}


func GetCounter(id int) (Wod, error) {
	for _, item := range itemCollection["Counter"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetFridge(id int) (Wod, error) {
	for _, item := range itemCollection["Fridge"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetTable(id int) (Wod, error) {
	for _, item := range itemCollection["Table"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDeco(id int) (Wod, error) {
	for _, item := range itemCollection["Deco"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetChair(id int) (Wod, error) {
	for _, item := range itemCollection["Chair"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetWallobject(id int) (Wod, error) {
	for _, item := range itemCollection["Wallobject"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetDish(id int) (Wod, error) {
	for _, item := range itemCollection["Dish"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetVendingmachine(id int) (Wod, error) {
	for _, item := range itemCollection["Vendingmachine"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetExpansion(id int) (Wod, error) {
	for _, item := range itemCollection["Expansion"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetAchievement(id int) (Wod, error) {
	for _, item := range itemCollection["Achievement"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetCoop(id int) (Wod, error) {
	for _, item := range itemCollection["Coop"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

func GetFastfood(id int) (Wod, error) {
	for _, item := range itemCollection["Fastfood"] {
		if id == item.ID {
			return item, nil
		}
	}
  return Wod{}, fmt.Errorf("No item found with id: %v", id)
}

