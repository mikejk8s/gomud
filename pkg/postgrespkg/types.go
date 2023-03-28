package postgrespkg

import (
	"database/sql"
	"os"

	"gorm.io/gorm"
)

var (
	POSTGRESS_HOST         = os.Getenv("POSTGRES_HOST")
	POSTGRES_PORT          = os.Getenv("POSTGRES_PORT")
	POSTGRES_USER          = os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD      = os.Getenv("POSTGRES_PASSWORD")
	POSTGRES_SSLMODE       = os.Getenv("POSTGRES_SSLMODE")
	POSTGRES_USERS_DB      = os.Getenv("POSTGRES_USERS_DB")      // users
	POSTGRES_CHARACTERS_DB = os.Getenv("POSTGRES_CHARACTERS_DB") // characters
	RUNNING_ON_DOCKER      = os.Getenv("RUNNING_ON_DOCKER")
)

type SqlConn struct {
	DB    *gorm.DB
	SqlDB *sql.DB // for closing the connection, and other things specific only to sql.DB
}
