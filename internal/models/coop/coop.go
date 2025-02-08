package coop

import (
	"cafego/internal/models/simple"
	"time"
)

type Coop struct {
	ID      int             `gorm:"primaryKey;autoIncrement;column:id"`
	Host    int             `gorm:"column:host;not null"`
	Members simple.IntSlice `gorm:"column:members;type:text;not null"`
	Kind    int             `gorm:"column:type;not null"` // Note: 'type' in SQL, 'Kind' in struct
	Dishes  simple.IntMap   `gorm:"column:dishes;type:text;not null"`
	Start   time.Time       `gorm:"column:start;default:CURRENT_TIMESTAMP"`
	End     time.Time       `gorm:"column:end;not null"`
}

func (coop Coop) TableName() string {
	return "coop"
}
