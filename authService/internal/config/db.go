package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		DBConfig.User, DBConfig.Password, DBConfig.DBName, DBConfig.SSLMode)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}

func CloseDB() {
	err := DB.Close()
	if err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}

	fmt.Println("Database connection closed!")
}
