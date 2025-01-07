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

var Wods []Wod
var Tiles []Wod
var Walls []Wod
var Doors []Wod
var Stoves []Wod
var Counters []Wod
var Fridges []Wod
var Tables []Wod
var Decors []Wod
var Chairs []Wod
var Wallobjects []Wod
var Dishes []Wod
var Ingredients []Wod
var Vendingmachines []Wod
var Expansions []Wod
var Achievements []Wod
var Coops []Wod
var FastFoods []Wod

func ReadAndCacheItems() {
	xmlFile, err := os.Open("./data/CafeItems.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()
	fmt.Println("Successfully Opened CafeItems.xml")

	var result Data
	decoder := xml.NewDecoder(xmlFile)
	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	Wods = result.Wods
	Tiles = []Wod{}
	Walls = []Wod{}
	Doors = []Wod{}
	Stoves = []Wod{}
	Counters = []Wod{}
	Fridges = []Wod{}
	Tables = []Wod{}
	Decors = []Wod{}
	Chairs = []Wod{}
	Wallobjects = []Wod{}
	Dishes = []Wod{}
	Ingredients = []Wod{}
	Vendingmachines = []Wod{}
	Expansions = []Wod{}
	Achievements = []Wod{}
	Coops = []Wod{}
	FastFoods = []Wod{}

	for _, wod := range Wods {
		switch wod.Group {
		case "Tile":
			Tiles = append(Tiles, wod)
		case "Wall":
			Walls = append(Walls, wod)
		case "Door":
			Doors = append(Doors, wod)
		case "Stove":
			Stoves = append(Stoves, wod)
		case "Counter":
			Counters = append(Counters, wod)
		case "Fridge":
			Fridges = append(Fridges, wod)
		case "Table":
			Tables = append(Tables, wod)
		case "Deco":
			Decors = append(Decors, wod)
		case "Chair":
			Chairs = append(Chairs, wod)
		case "Wallobject":
			Wallobjects = append(Wallobjects, wod)
		case "Dish":
			Dishes = append(Dishes, wod)
		case "Ingredient":
			Ingredients = append(Ingredients, wod)
		case "Vendingmachine":
			Vendingmachines = append(Vendingmachines, wod)
		case "Expansion":
			Expansions = append(Expansions, wod)
		case "Achievement":
			Achievements = append(Achievements, wod)
		case "Coop":
			Coops = append(Coops, wod)
		case "Fastfood":
			FastFoods = append(FastFoods, wod)
		}
	}

	fmt.Printf("Successfully loaded %d WOD entries\n", len(Wods))
	fmt.Println("Filtered and grouped WOD entries")
}
