package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func InitDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it : ", err)
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL not found in the .env file ")
	}
	fmt.Println("url : ", dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Failed to open a connection to the database ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Database is unreachable ", err)
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute)
	fmt.Println("✅ Connected to PostgreSQL")
	return db
}

func GetENVdata(key string) any {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it : ", err)
	}
	var value any
	if key == "JWT_SECRET" {
		value = []byte(os.Getenv(key))
	} else {
		value = os.Getenv(key)
	}
	if value == "" {
		log.Println("%w Not found ", key)
	}
	return value
}
