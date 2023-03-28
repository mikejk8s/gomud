package models

import (
	"time"
)

type Character struct {
	ID             uint64 `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"uniqueIndex" json:"name"`
	Class          string
	Race           string `gorm:"default:'HUMAN'" json:"race"`
	Level          int    `gorm:"default:1" json:"level"`
	CreatedAt      time.Time
	Alive          bool   `gorm:"default:true" json:"alive"`
	CharacterOwner string `gorm:"default:'player'" json:"character_owner"`
}
