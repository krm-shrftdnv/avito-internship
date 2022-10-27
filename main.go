package main

import (
	"avito-internship/app"
	"avito-internship/handler"
	"fmt"
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	app.Init()

	app.FiberApp.Get("/balance/:userId", handler.GetBalance)
	app.FiberApp.Post("/balance/add", handler.AddBalance)
	app.FiberApp.Post("/reserve", handler.Reserve)
	app.FiberApp.Post("/reserve/approve", handler.ApproveReserve)
	app.FiberApp.Get("/report/create", handler.GetReport)
	app.FiberApp.Get("/report/:filename", handler.DownloadReport)
	app.FiberApp.Get("/balance/:userId/operations", handler.GetOperations)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatalln(app.FiberApp.Listen(fmt.Sprintf(":%v", port)))
}
