package postgrespkg

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dbName string) (*SqlConn, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		POSTGRESS_HOST, POSTGRES_USER, POSTGRES_PASSWORD,
		dbName, POSTGRES_PORT, POSTGRES_SSLMODE,
		os.Getenv("TZ"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	var sqlUtil = new(SqlConn)
	sqlUtil = &SqlConn{DB: db, SqlDB: sqlDB}
	return sqlUtil, nil
}

func (s *SqlConn) ConnectSQLToSchema(schemaName string) error {
	db := s.DB.Exec(fmt.Sprintf("SET search_path TO %s", schemaName))
	if db.Error != nil {
		return db.Error
	}
	return nil
}

func (s *SqlConn) Close() error {
	return s.SqlDB.Close()
}
