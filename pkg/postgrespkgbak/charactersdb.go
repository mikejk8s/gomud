package postgrespkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/mikejk8s/gmud/pkg/models"
)

func CreateCharacterTable(db *gorm.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.WithContext(ctx).AutoMigrate(&models.Character{})
	if err != nil {
		return fmt.Errorf("error migrating Character table: %w", err)
	}

	log.Printf("Character table created successfully.")
	return nil
}
