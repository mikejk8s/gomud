package postgrespkg

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LoginReq struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type SqlConn struct {
	DB *gorm.DB
}

func (conn *SqlConn) GetSQLConn(dbname string) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s", Hostname, Username, Password, "postgres", SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error:", err.Error())
		return err
	}
	conn.DB = db

	if err := conn.DB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", dbname)).Error; err != nil {
		return fmt.Errorf("error creating Character database: %w", err)
	}

	return nil
}

var RunningOnDocker = false

var (
	Username = "gmud"
	Password = "gmud"
	Hostname = "127.0.0.1"
	SSLMode  = "disable"
)
