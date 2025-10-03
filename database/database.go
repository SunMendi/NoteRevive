package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {

	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .Env File")
	}

	dbname := os.Getenv("DB_NAME")
	dbuser := os.Getenv("DB_USER")
	dbport := os.Getenv("DB_PORT")
	dbhost := os.Getenv("DB_HOST")
	dbpass := os.Getenv("DB_PASS")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, dbname)

	if dsn == "" {
		log.Panic("Error Loading .Env File")
	}

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// log.Printf("DB_HOST: %s", dbhost)
	// log.Printf("DB_USER: %s", dbuser)
	// log.Printf("DB_PASS: %s", dbpass)
	// log.Printf("DSN: %s", dsn)

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
		return nil, err
	}

	return DB, nil
}
