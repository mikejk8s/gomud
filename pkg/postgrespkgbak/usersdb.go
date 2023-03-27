package postgrespkg

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mikejk8s/gmud/pkg/models"
)

func CreateNewUser(db *gorm.DB, userInfo models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &models.User{
		CreatedAt:    time.Now(),
		Name:         userInfo.Name,
		Username:     userInfo.Username,
		Email:        userInfo.Email,
		PasswordHash: string(hashedPassword),
		RememberHash: uuid.New().String(),
	}

	if err := db.Create(user).Error; err != nil {
		return err
	}

	return nil
}

func CreateUsersTable(db *gorm.DB) error {
	err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
        id              SERIAL PRIMARY KEY,
        created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP,
        deleted_at      TIMESTAMP,
        name            VARCHAR(255),
        username        VARCHAR(255) UNIQUE NOT NULL,
        email           VARCHAR(255) UNIQUE NOT NULL,
        password_hash   VARCHAR(255),
        remember_hash   VARCHAR(255)
    );`).Error

	if err != nil {
		return err
	}

	return nil
}

func UsersMigration(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	return nil
}
