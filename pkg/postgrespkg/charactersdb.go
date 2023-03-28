package postgrespkg

import (
	"github.com/mikejk8s/gmud/pkg/models"
	"gorm.io/gorm"
)

func CreateCharacterTable(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Character{})
	if err != nil {
		return err
	}
	return nil
}

func (s *SqlConn) CreateNewCharacter(character *models.Character) error {
	if err := s.DB.Create(character).Error; err != nil {
		return err
	}
	return nil
}
