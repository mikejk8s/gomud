package postgrespkg

import (
	"fmt"
	"os"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Conn struct {
	DB *gorm.DB
}

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_SSLMODE"),
		os.Getenv("POSTGRES_TIMEZONE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (conn *Conn) GetSQLConn(schemaName string) error {
	db := conn.DB.Exec(fmt.Sprintf("SET search_path TO %s", schemaName))
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func Close(db *gorm.DB) {
	dbConn, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	dbConn.Close()
}