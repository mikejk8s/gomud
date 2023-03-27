package postgrespkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func (conn *SqlConn) CreateCharacterTable2() error {
	db, err := gorm.Open(postgres.Open(conn.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("error opening connection to the database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.WithContext(ctx).AutoMigrate(&Character{})
	if err != nil {
		return fmt.Errorf("error migrating Character table: %w", err)
	}

	log.Printf("Character table created successfully.")
	return nil
}