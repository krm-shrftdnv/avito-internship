package app

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DataBase *sql.DB
var FiberApp *fiber.App
var Validate *validator.Validate

type Error struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Fields  interface{} `json:"fields"`
}

type ResponseBody struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func env() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

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

func fiberInit() {
	FiberApp = fiber.New(fiber.Config{
		AppName: "avito-internship",
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			err = ctx.Status(code).JSON(ResponseBody{Data: nil, Error: Error{Message: err.Error(), Code: code}})
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			return nil
		},
	})
	FiberApp.Use(recover.New())
}

func validatorInit() {
	Validate = validator.New()
}

func Init() {
	env()
	db()
	fiberInit()
	validatorInit()
}
