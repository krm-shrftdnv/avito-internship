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
		err := recover()
		if err != nil {
			log.Println(err)
		}
	}()

	app.Init()

	app.FiberApp.Get("/balance", handler.GetBalance)
	app.FiberApp.Post("/balance/add", handler.AddBalance)
	app.FiberApp.Post("/reserve", handler.Reserve)
	app.FiberApp.Post("/reserve/approve", handler.ApproveReserve)
	app.FiberApp.Get("/report", handler.GetReport)
	app.FiberApp.Get("/transactions", handler.GetTransactions)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000" // docker
		//port = "3001" // todo: not push debug
	}
	log.Fatalln(app.FiberApp.Listen(fmt.Sprintf(":%v", port)))
}
