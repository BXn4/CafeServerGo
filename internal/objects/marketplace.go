package objects

import (
	"fmt"
)

func NewMarketplace() (*Cafe, error) {

	cafe := &Cafe{
		id:          -1,
		playerID:    -1,
		playerStart: [2]int{1, 2},
		ownerName:   "Marketplace",
		background:  MarketplaceBackground,
		rating:      0,
		luxury:      -1,
		size:        12,
	}

	err := cafe.ParseTiles("7+7+7+7+7+7+7+7+7+7+7+7+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3")

	if err != nil {
		return nil, fmt.Errorf("Can not parse tiles for marketplace: %v\n", err)
	}

	err = cafe.ParseObjects("6+6+511+0#5+11+515+0#1+9+512+0#11+12+513+0#12+11+514+0")
	if err != nil {
		return nil, fmt.Errorf("Can not parse objects for marketplace: %v\n", err)
	}

	return cafe, nil
}
