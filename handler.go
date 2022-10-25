package main

import (
	"avito-internship/app"
	"avito-internship/models"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"strconv"
)

func GetBalance(c *fiber.Ctx) (err error) {
	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	id, err := strconv.Atoi(c.Query("id"))
	userModel := models.User{Id: int32(id)}
	user, err := userModel.GetById(transaction)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.JSON(user)
}

func AddBalance(c *fiber.Ctx) (err error) {
	return nil
}

func Reserve(c *fiber.Ctx) (err error) {
	return nil
}

func ApproveReserve(c *fiber.Ctx) (err error) {
	return nil
}

func GetReport(c *fiber.Ctx) (err error) {
	return nil
}

func GetTransactions(c *fiber.Ctx) (err error) {
	return nil
}
