package main

import (
	"avito-internship/app"
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

	app.Init()

	fiberApp := fiber.New()
	fiberApp.Get("/balance", GetBalance)
	fiberApp.Post("/balance/add", AddBalance)
	fiberApp.Post("/reserve", Reserve)
	fiberApp.Post("/reserve/approve", ApproveReserve)
	fiberApp.Get("/report", GetReport)
	fiberApp.Get("/transactions", GetTransactions)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(fiberApp.Listen(fmt.Sprintf(":%v", port)))
}
