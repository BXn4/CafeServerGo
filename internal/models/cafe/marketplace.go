package cafe

import (
	"cafego/internal/models/object"
	"cafego/internal/models/simple"
	"fmt"
)

func NewMarketplace(id int) (*Cafe, error) {
	// Check if id is negative
	// Cafe is now having room types, so we can sepperate it
	/* if id >= 0 {
	return nil, fmt.Errorf("Marketplace id cant be possitive!")
	} */

	f := &simple.IntMatrix{}
	tileString := "7+7+7+7+7+7+7+7+7+7+7+7+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+7+7+7+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3+7+3+3+3+3+3+3+3+3+3+3+3"
	err := f.Scan([]byte(tileString))
	if err != nil {
		return nil, fmt.Errorf("Can not parse tiles for marketplace: %v", err)
	}

	objs := &object.ObjectList{}
	objsStr := "6+6+511+0#5+11+515+0#1+9+512+0#11+12+513+0#12+11+514+0"
	err = objs.Scan([]byte(objsStr))
	if err != nil {
		return nil, fmt.Errorf("Can not parse objects for marketplace: %v", err)
	}

	return &Cafe{
		ID:         id,
		PlayerID:   id,
		OwnerName:  "Marketplace",
		Background: MarketplaceBackground,
		Rating:     0,
		Luxury:     -1,
		Size:       12,
		Objects:    *objs,
		Tiles:      *f,
		roomType:   MarketRoom,
	}, nil
}
