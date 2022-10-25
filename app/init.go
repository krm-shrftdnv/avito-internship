package app

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DataBase *sql.DB

func db() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"database", 5432, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), "balance_db")
	var err error
	DataBase, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed initialize db. %s", err.Error()))
	}
	err = DataBase.Ping()
	if err != nil {
		log.Println(fmt.Errorf("failed ping db. %s", err.Error()))
	}
	log.Println("Database connection established")
}

func env() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func Init() {
	env()
	db()
}
