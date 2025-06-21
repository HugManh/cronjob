package db

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(dbhost string, dbport int, dbuser, dbpass, dbname string, usessl bool) (*gorm.DB, error) {
	dbssl := "disable"
	if usessl {
		dbssl = "require"
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbhost, dbport, dbuser, dbpass, dbname, dbssl,
	)
	log.Printf("Connecting to database with DSN: %s", dsn)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
