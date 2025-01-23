package daos

import (
	"cafego/internal/types/cafetypes"
	"time"
)

type CafeObjectDAO struct {
	Kind     cafetypes.CafeObjectKind     `json:"id"`
	Pos      [2]int                       `json:"pos"`
	Rotation cafetypes.CafeObjectRotation `json:"rotation"`

	DishID     int `json:"dish_id,omitempty"`
	DishStatus int `json:"dish_status,omitempty"`
	DishAmount int `json:"dish_amount,omitempty"`

	FancyIng   bool       `json:"fancy_ing,omitempty"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishesAt *time.Time `json:"finishes_at,omitempty"`
}
