package postgrespkg

import (
	"time"

	"gorm.io/gorm"
)

type Character struct {
	gorm.Model
	Name           string `gorm:"not null;uniqueIndex"`
	Class          string `gorm:"not null"`
	Race           string `gorm:"not null;default:'HUMAN'"`
	Level          int    `gorm:"not null;default:1"`
	CreatedAt      time.Time
	Alive          bool   `gorm:"not null;default:true"`
	Characterowner string `gorm:"not null;default:'player'"`
}

func CreateCharacterTable(db *gorm.DB) error {
	err := db.AutoMigrate(&Character{})
	if err != nil {
		return err
	}
	return nil
}

func CreateNewCharacter(db *gorm.DB, name string, class string, race string, level int, alive bool) error {
	character := &Character{
		Name:      name,
		Class:     class,
		Race:      race,
		Level:     level,
		CreatedAt: time.Now(),
		Alive:     alive,
	}

	if err := db.Create(character).Error; err != nil {
		return err
	}
	return nil
}