package postgrespkg

import (
	"log"
	"time"

	"github.com/mikejk8s/gmud/logger"
	"github.com/mikejk8s/gmud/pkg/models"
)

func (s *SqlConn) CreateNewUser(name string, password string) error {
	user := &models.User{
		CreatedAt:    time.Now(),
		Name:         name,
		PasswordHash: password,
		RememberHash: "",
	}

	if err := s.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetCharacterList returns an array of characters associated with the account accOwner.
func (s *SqlConn) GetCharacterList(accOwner string) ([]*models.Character, error) {
	cDBLogger := logger.GetNewLogger()
	err := cDBLogger.AssignOutput("characterDB", "./logs/characterDBconn")
	if err != nil {
		log.Println(err)
	}
	if err != nil {
		cDBLogger.LogUtil.Errorf("Error %s connecting to characterDB during fetching the %s accounts characters: ", err, accOwner)
		panic(err.Error())
	}
	return s.GetCharactersByUserName(accOwner)
}

func (s *SqlConn) GetCharactersByUserName(name string) ([]*models.Character, error) {
	var characters []*models.Character
	if err := s.DB.Where("characterowner = ?", name).Find(&characters).Error; err != nil {
		return nil, err
	}
	return characters, nil
}
