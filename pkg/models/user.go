package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint64 `gorm:"primaryKey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string
	PasswordHash string
	RememberHash string
}

// CheckPassword gets plain password as input and checks if it matches the hashed password in the database.
//
// user.Password is set during fetching the users database, and retrieved as already hashed.
//
// If err is not nil, then the password is not correct and SSH password authentication will fail by returning false.
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(providedPassword))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
