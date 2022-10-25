package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DataBase *sql.Conn

func db() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"database", 5432, os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), "balance_db")
	DataBase, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(fmt.Errorf("failed initialize db. %s", err.Error()))
	}
	err = DataBase.Ping()
	if err != nil {
		log.Println(fmt.Errorf("failed ping db. %s", err.Error()))
	}
	log.Println("Database connection established")
}

func Init() {
	db()
}
