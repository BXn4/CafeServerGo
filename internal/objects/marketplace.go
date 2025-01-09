package objects

import (
  "fmt"
)


func NewMarketplace() (*Cafe, error) {


  cafe := &Cafe{
    ID: -1,
    PlayerID: -1,
    PlayerStart: []int{1,2},
    OwnerName: "Marketplace",
    Background: MarketplaceBackground,
    Rating: 0,
    Luxury: -1,
    ExpansionID: 4,
    Size: 12,
  } 

  err := cafe.ParseTiles("7+7+7+7+7+7+7+7+7+7+7+7+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3")

  if err != nil {
    fmt.Printf("Can not parse tiles for marketplace!\n")
    return nil, err
  }

  err = cafe.ParseObjects("6+6+511+0#5+11+515+0#1+9+512+0#11+12+513+0#12+11+514+0")
  if err != nil {
    fmt.Printf("Can not parse objects for marketplace!\n")
    return nil, err
  }

  return cafe, nil 
}

