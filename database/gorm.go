package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm" // golang object-relational mapping
)

// Singleton Design pattern for database connection

var (
	dbInstance *gorm.DB
	models     []any
	once       sync.Once
)

func RegisterModels(m ...any) {
	models = append(models, m...)
}

func DB() (*gorm.DB, error) {
	var err error
	initiateDB := func() {
		maxRetries := 5
		retryInterval := time.Second * 5

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		for i := 0; i < maxRetries; i++ {
			dbInstance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err == nil {
				break
			}
			log.Printf("failed to connect to database (attempt %d/%d): %v\n", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(retryInterval)
			}
		}
		if err != nil {
			log.Printf("failed to connect to database after %d attempts: %v\n", maxRetries, err)
			return
		}

		err := dbInstance.AutoMigrate(models...)
		if err != nil {
			log.Printf("failed to migrate database: %v\n", err)
			return
		}
	}

	once.Do(initiateDB)
	return dbInstance, err
}
