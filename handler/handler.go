package handler

import (
	"avito-internship/app"
	"avito-internship/models"
	"database/sql"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"log"
	"strconv"
)

type ReserveForm struct {
	UserId    int     `json:"id" validate:"required,numeric,min=1"`
	ServiceId int     `json:"service_id" validate:"required,numeric,min=1"`
	OrderId   int     `json:"order_id" validate:"required,numeric,min=1"`
	Price     float32 `json:"price" validate:"required,number"`
}

type ReserveApproveForm struct {
	UserId    int     `json:"id" validate:"required,numeric,min=1"`
	ServiceId int     `json:"service_id" validate:"required,numeric,min=1"`
	OrderId   int     `json:"order_id" validate:"required,numeric,min=1"`
	Sum       float32 `json:"sum" validate:"required,number"`
}

type ValidationError struct {
	FailedField string
	Tag         string
	Value       string
}

func validate(s interface{}) []*ValidationError {
	var errs []*ValidationError
	err := app.Validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errs = append(errs, &element)
		}
	}
	return errs
}

func GetBalance(c *fiber.Ctx) (err error) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User Id is required")
	}
	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
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
	errs := validate(user)
	if errs != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(app.ResponseBody{Error: app.Error{Code: fiber.StatusUnprocessableEntity, Fields: errs}})
	}

	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	userModel := &models.User{Id: user.Id}
	userModel, err = userModel.GetById(transaction)
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
	// todo: add operation
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
	errs := validate(form)
	if errs != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(app.ResponseBody{Error: app.Error{Code: fiber.StatusUnprocessableEntity, Fields: errs}})
	}

	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	user := &models.User{Id: int32(form.UserId)}
	user, err = user.GetById(transaction)
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
	if user.BalanceValue < form.Price {
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
	reserve := models.Reserve{UserId: user.Id, OrderId: int32(form.OrderId), ServiceId: int32(form.ServiceId), Value: form.Price}
	err = reserve.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error while reserving. Try again")
	}
	user.BalanceValue -= form.Price
	err = user.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	// todo: add operation
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.JSON(app.ResponseBody{Data: reserve})
}

func ApproveReserve(c *fiber.Ctx) (err error) {
	form := new(ReserveApproveForm)
	if err = c.BodyParser(form); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	errs := validate(form)
	if errs != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(app.ResponseBody{Error: app.Error{Code: fiber.StatusUnprocessableEntity, Fields: errs}})
	}

	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	reserve := &models.Reserve{OrderId: int32(form.OrderId), ServiceId: int32(form.ServiceId)}
	reserve, err = reserve.GetByOrderIdAndServiceId(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		if err == sql.ErrNoRows {
			return fiber.NewError(fiber.StatusNotFound, "Reserve not found")
		}
		return errors.Wrap(err, "")
	}
	if reserve.IsDebited {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return fiber.NewError(fiber.StatusBadRequest, "Reserve is already approved")
	}
	if reserve.Value != form.Sum {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return fiber.NewError(fiber.StatusBadRequest, "Reserved sum not equals entered sum")
	}
	reserve.IsDebited = true
	err = reserve.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}
		log.Println(err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error while updating reserve. Try again")
	}
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}
	return c.JSON(app.ResponseBody{Data: reserve})
}

func GetReport(c *fiber.Ctx) (err error) {
	return nil
}

func GetTransactions(c *fiber.Ctx) (err error) {
	return nil
}
