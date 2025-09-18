package db

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var global *gorm.DB

// Open initializes the global DB connection using the provided DSN.
func Open(dsn string) (*gorm.DB, error) {
	if global != nil {
		return global, nil
	}
	gormLogger := logger.New(
		log.New(log.Writer(), "gorm: ", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, err
	}
	global = db
	return db, nil
}

// DB returns the initialized global database connection.
func DB() *gorm.DB {
	return global
}
