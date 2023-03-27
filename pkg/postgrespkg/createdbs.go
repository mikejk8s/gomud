package postgrespkg

import (
	"fmt"
	"log"
)

// CreateDatabases creates two databases with the names specified by env vars.
func CreateDatabases() error {
	databases := []string{POSTGRES_DB1, POSTGRES_DB2}

	for _, dbName := range databases {
		db, err := Connect()
		if err != nil {
			return fmt.Errorf("error connecting to database: %w", err)
		}

		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return fmt.Errorf("error creating database %s: %w", dbName, err)
		}
		log.Printf("Database %s created successfully.", dbName)

		err = db.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, POSTGRES_USER)).Error
		if err != nil {
			return fmt.Errorf("error granting privileges to database %s: %w", dbName, err)
		}
		log.Printf("Privileges granted to user %s for database %s.", POSTGRES_USER, dbName)

		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("error getting sql.DB from gorm.DB: %w", err)
		}
		err = sqlDB.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
	}

	return nil
}
