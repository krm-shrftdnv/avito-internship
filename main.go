package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	defer func() {
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	Init()

	app := fiber.New()
	app.Get("/balance", GetBalance)
	app.Post("/balance/add", AddBalance)
	app.Post("/reserve", Reserve)
	app.Post("/reserve/approve", ApproveReserve)
	app.Get("/report", GetReport)
	app.Get("/transactions", GetTransactions)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
