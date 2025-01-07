package utils

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Wod struct {
	ID    int    `xml:"id,attr"`
	Name  string `xml:"n,attr"`
	Group string `xml:"g,attr"`
	Type  string `xml:"t,attr"`
}

type Tile struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Wall struct {
	Wod
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Door struct {
	Wod
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Stove struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Counter struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Fridge struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
	InventorySize int `xml:"inventorySize,omitempty"`
}

type Table struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Decor struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Chair struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Wallobject struct {
	Wod
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Dish struct {
	Wod
	Events           int    `xml:"events,omitempty"`
	XP               int    `xml:"xp,omitempty"`
	IncomePerServing int    `xml:"incomePerServing,omitempty"`
	Servings         int    `xml:"servings,omitempty"`
	Duration         int    `xml:"Duration,omitempty"`
	Level            int    `xml:"level,omitempty"`
	Category         int    `xml:"dishcategory,omitempty"`
	Requirements     string `xml:"requirements,omitempty"`
}

type Ingredient struct {
	Wod
	Cash          int    `xml:"cash,omitempty"`
	Gold          int    `xml:"gold,omitempty"`
	Friends       int    `xml:"friends,omitempty"`
	GoldNoFriends int    `xml:"goldNoFriends,omitempty"`
	Events        int    `xml:"events,omitempty"`
	XP            int    `xml:"xp,omitempty"`
	Level         int    `xml:"level,omitempty"`
	Category      string `xml:"category,omitempty"`
	Amount        int    `xml:"amount,omitempty"`
}

type Vendingmachine struct {
	Wod
	Walkable      int `xml:"walkable,attr,omitempty"`
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
}

type Expansion struct {
	Wod
	Cash          int `xml:"cash,omitempty"`
	Gold          int `xml:"gold,omitempty"`
	Friends       int `xml:"friends,omitempty"`
	GoldNoFriends int `xml:"goldNoFriends,omitempty"`
	Events        int `xml:"events,omitempty"`
	XP            int `xml:"xp,omitempty"`
	Level         int `xml:"level,omitempty"`
	ExpansionID   int `xml:"expansionID,attr,omitempty"`
	SizeX         int `xml:"sizeX,attr,omitempty"`
	SizeY         int `xml:"sizeY,attr,omitempty"`
}

type Achievement struct {
	Wod
}

type Coop struct {
	Wod
	Events     int    `xml:"events,omitempty"`
	MaxMembers int    `xml:"maxMember,omitempty"`
	MaxLevel   int    `xml:"maxLevel,omitempty"`
	Chips      int    `xml:"chips,omitempty"`
	XP         int    `xml:"xp,omitempty"`
	Gold       int    `xml:"gold,omitempty"`
	Duration   int    `xml:"Duration,omitempty"`
	Dishes     string `xml:"dishes,omitempty"`
}

type FastFood struct {
	Wod
	Cash             int `xml:"cash,omitempty"`
	Gold             int `xml:"gold,omitempty"`
	Friends          int `xml:"friends,omitempty"`
	GoldNoFriends    int `xml:"goldNoFriends,omitempty"`
	Events           int `xml:"events,omitempty"`
	XP               int `xml:"xp,omitempty"`
	Level            int `xml:"level,omitempty"`
	IncomePerServing int `xml:"incomePerServing,omitempty"`
	Servings         int `xml:"servings,omitempty"`
	RatingBonus      int `xml:"ratingBonus,omitempty"`
}

type Data struct {
	Wods []Wod `xml:"wod"`
}

var Wods []Wod
var Tiles []Tile
var Walls []Wall
var Doors []Door
var Stoves []Stove
var Counters []Counter
var Fridges []Fridge
var Tables []Table
var Decors []Decor
var Chairs []Chair
var Wallobjects []Wallobject
var Dishes []Dish
var Ingredients []Ingredient
var Vendingmachines []Vendingmachine
var Expansions []Expansion
var Achievements []Achievement
var Coops []Coop
var FastFoods []FastFood

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
	fmt.Printf("Successfully loaded %d WOD entries\n", len(Wods))

	FilterWodsByGroup()
}

func FilterWodsByGroup() {
	Tiles = []Tile{}
	Walls = []Wall{}
	Doors = []Door{}
	Stoves = []Stove{}
	Counters = []Counter{}
	Fridges = []Fridge{}
	Tables = []Table{}
	Decors = []Decor{}
	Chairs = []Chair{}
	Wallobjects = []Wallobject{}
	Dishes = []Dish{}
	Ingredients = []Ingredient{}
	Vendingmachines = []Vendingmachine{}
	Expansions = []Expansion{}
	Achievements = []Achievement{}
	Coops = []Coop{}
	FastFoods = []FastFood{}

	for _, wod := range Wods {
		switch wod.Group {
		case "Tile":
			tile := Tile{Wod: wod}
			Tiles = append(Tiles, tile)
		case "Wall":
			wall := Wall{Wod: wod}
			Walls = append(Walls, wall)
		case "Door":
			door := Door{Wod: wod}
			Doors = append(Doors, door)
		case "Stove":
			stove := Stove{Wod: wod}
			Stoves = append(Stoves, stove)
		case "Counter":
			counter := Counter{Wod: wod}
			Counters = append(Counters, counter)
		case "Fridge":
			fridge := Fridge{Wod: wod}
			Fridges = append(Fridges, fridge)
		case "Table":
			table := Table{Wod: wod}
			Tables = append(Tables, table)
		case "Decor":
			decor := Decor{Wod: wod}
			Decors = append(Decors, decor)
		case "Chair":
			chair := Chair{Wod: wod}
			Chairs = append(Chairs, chair)
		case "Wallobject":
			wallobject := Wallobject{Wod: wod}
			Wallobjects = append(Wallobjects, wallobject)
		case "Dish":
			dish := Dish{Wod: wod}
			Dishes = append(Dishes, dish)
		case "Ingredient":
			ingredient := Ingredient{Wod: wod}
			Ingredients = append(Ingredients, ingredient)
		case "Vendingmachine":
			vendingmachine := Vendingmachine{Wod: wod}
			Vendingmachines = append(Vendingmachines, vendingmachine)
		case "Expansion":
			expansion := Expansion{Wod: wod}
			Expansions = append(Expansions, expansion)
		case "Achievement":
			achievement := Achievement{Wod: wod}
			Achievements = append(Achievements, achievement)
		case "Coop":
			coop := Coop{Wod: wod}
			Coops = append(Coops, coop)
		case "FastFood":
			fastFood := FastFood{Wod: wod}
			FastFoods = append(FastFoods, fastFood)
		}
	}

	fmt.Println("Successfully filtered all WOD into groups.")
}
