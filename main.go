package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	app := fiber.New()

	app.Get("/balance", getBalance)
	app.Post("/balance/add", addBalance)
	app.Post("/reserve", reserve)
	app.Post("/reserve/approve", approveReserve)
	app.Get("/report", getReport)
	app.Get("/transactions", getTransactions)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}

func getBalance(c *fiber.Ctx) (err error) {
	return nil
}

func addBalance(c *fiber.Ctx) (err error) {
	return nil
}

func reserve(c *fiber.Ctx) (err error) {
	return nil
}

func approveReserve(c *fiber.Ctx) (err error) {
	return nil
}

func getReport(c *fiber.Ctx) (err error) {
	return nil
}

func getTransactions(c *fiber.Ctx) (err error) {
	return nil
}
