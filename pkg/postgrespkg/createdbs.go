package postgrespkg

import (
	"fmt"
	"log"

	"github.com/mikejk8s/gmud/pkg/models"
)

// CreateDatabases creates two databases with the names specified by env vars.
func (s *SqlConn) CreateDatabases() error {
	databases := []string{POSTGRES_USERS_DB, POSTGRES_CHARACTERS_DB}

	for _, dbName := range databases {
		err := s.DB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return fmt.Errorf("error creating database %s: %w", dbName, err)
		}
		log.Printf("Database %s created successfully.", dbName)

		err = s.DB.Exec(fmt.Sprintf("GRANT ALL PRIVILEGES ON DATABASE %s TO %s", dbName, POSTGRES_USER)).Error
		if err != nil {
			return fmt.Errorf("error granting privileges to database %s: %w", dbName, err)
		}
		log.Printf("Privileges granted to user %s for database %s.", POSTGRES_USER, dbName)
		err = s.SqlDB.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %w", err)
		}
	}

	return nil
}

func (s *SqlConn) CreateUsersTable() error {
	return s.DB.Exec(`CREATE TABLE IF NOT EXISTS users (
        id              SERIAL PRIMARY KEY,
        created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at      TIMESTAMP,
        deleted_at      TIMESTAMP,
        name            VARCHAR(255),
        password_hash   VARCHAR(255),
        remember_hash   VARCHAR(255)
    );`).Error
}

func (s *SqlConn) CreateCharacterTable() error {
	return s.DB.Exec(`CREATE TABLE IF NOT EXISTS characters (
		id              BIGSERIAL PRIMARY KEY,
		name            VARCHAR(30) UNIQUE NOT NULL,
		class           VARCHAR(15) NOT NULL,
		race            VARCHAR(15) NOT NULL DEFAULT 'HUMAN',
		level           INTEGER NOT NULL DEFAULT '1',
		created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		alive           BOOLEAN NOT NULL DEFAULT '1',
		characterowner  VARCHAR(20) NOT NULL DEFAULT 'player'
	);`).Error
}

func (s *SqlConn) MigrateCharacters() error {
	return s.DB.AutoMigrate(&models.Character{})
}

func (s *SqlConn) MigrateUsers() error {
	return s.DB.AutoMigrate(&models.User{})
}
