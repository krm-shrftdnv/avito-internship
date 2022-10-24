package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
)

func main() {
	app := fiber.New()

	//app.Get("/", indexHandler) // Add this
	//
	//app.Post("/", postHandler) // Add this
	//
	//app.Put("/update", putHandler) // Add this
	//
	//app.Delete("/delete", deleteHandler) // Add this

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
}
