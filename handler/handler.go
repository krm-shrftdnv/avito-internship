package handler

import (
	"avito-internship/app"
	"avito-internship/models"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"strconv"
)

type ReserveForm struct {
	UserId    int     `json:"id"`
	ServiceId int     `json:"service_id"`
	OrderId   int     `json:"order_id"`
	Price     float32 `json:"price"`
}

func GetBalance(c *fiber.Ctx) (err error) {
	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	id, err := strconv.Atoi(c.Query("id"))
	userModel := models.User{Id: int32(id)}
	user, err := userModel.GetById(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return errors.Wrap(err, "")
	}
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.JSON(app.ResponseBody{Data: user})
}

func AddBalance(c *fiber.Ctx) (err error) {
	user := new(models.User)
	if err = c.BodyParser(user); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	userModel, err := user.GetById(transaction)
	if err != nil {
		if err == sql.ErrNoRows {
			userModel = &models.User{
				Id:           user.Id,
				Name:         user.Name,
				BalanceValue: 0,
			}
		} else {
			tErr := transaction.Rollback()
			if tErr != nil {
				return errors.Wrap(tErr, "")
			}
			return errors.Wrap(err, "")
		}
	}
	userModel.BalanceValue += user.BalanceValue
	err = userModel.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.JSON(app.ResponseBody{Data: userModel})
}

func Reserve(c *fiber.Ctx) (err error) {
	form := new(ReserveForm)
	if err = c.BodyParser(form); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	user := models.User{Id: int32(form.UserId)}
	userModel, err := user.GetById(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return errors.Wrap(err, "")
	}
	if userModel.BalanceValue < form.Price {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}
		return fiber.NewError(fiber.StatusBadRequest, "Not enough money")
	}
	service := models.Service{Id: int32(form.ServiceId), Price: form.Price}
	err = service.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	order := models.Order{UserId: userModel.Id, ServiceId: service.Id}
	err = order.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	reserve := models.Reserve{UserId: userModel.Id, OrderId: order.Id, Value: form.Price}
	err = reserve.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	userModel.BalanceValue -= form.Price
	err = userModel.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}
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
