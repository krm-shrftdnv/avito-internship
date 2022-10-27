package handler

import (
	"avito-internship/app"
	"avito-internship/models"
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
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

type ReportParams struct {
	Year  int32 `json:"year" validation:"required,numeric,min=0,max=9999"`
	Month int32 `json:"month" validation:"required,numeric,min=1,max=12"`
}

type ReportRow struct {
	Sum  float32 `json:"sum"`
	Name string  `json:"name"`
}

type ReportUrl struct {
	Url string `json:"url"`
}

type OperationsFilter struct {
	Page    int    `validation:"numeric,min=1"`
	PerPage int    `validation:"numeric,min=1"`
	Sort    string `validation:"oneof=value created_at"`
	SortDir string `validation:"oneof=asc desc"`
}

func (reportRow *ReportRow) toString() []string {
	return []string{strconv.FormatFloat(float64(reportRow.Sum), 'f', 2, 64), reportRow.Name}
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
	id, err := strconv.Atoi(c.Params("userId"))
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
	err = userModel.Save(transaction)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	operation := models.Operation{
		UserId:     userModel.Id,
		IsIncrease: true,
		Value:      user.BalanceValue,
	}
	err = operation.Save(transaction)
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
	operation := models.Operation{
		UserId:     int32(form.UserId),
		IsIncrease: false,
		Value:      form.Sum,
	}
	err = operation.Save(transaction)
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
	return c.JSON(app.ResponseBody{Data: reserve})
}

func GetReport(c *fiber.Ctx) (err error) {
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		return errors.Wrap(err, "Incorrect year")
	}
	month, err := strconv.Atoi(c.Query("month"))
	if err != nil {
		return errors.Wrap(err, "Incorrect month")
	}
	reportParams := ReportParams{Year: int32(year), Month: int32(month)}
	errs := validate(reportParams)
	if errs != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(app.ResponseBody{Error: app.Error{Code: fiber.StatusUnprocessableEntity, Fields: errs}})
	}

	transaction, err := app.DataBase.Begin()
	if err != nil {
		return errors.Wrap(err, "")
	}
	rows, err := transaction.
		Query(
			"select sum(r.value) as sum, s.name as name from reserve r left join service s on r.service_id = s.id where created_at between ($1 || '-' || $2 || '-01 00:00:00')::timestamp and date_trunc('month', ($1 || '-' || $2 || '-01 00:00:00')::timestamp) + interval '1 month - 1 day' group by s.id;",
			reportParams.Year,
			reportParams.Month,
		)
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}

		return errors.Wrap(err, "")
	}
	directory := "./reports"
	fileName := fmt.Sprintf("%d-%d_services_report.csv", reportParams.Year, reportParams.Month)
	file, err := os.Create(filepath.Join(directory, fileName))
	if err != nil {
		err = os.Mkdir(directory, 0777)
		if err != nil {
			tErr := transaction.Rollback()
			if tErr != nil {
				return errors.Wrap(tErr, "")
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Error while forming report")
		} else {
			file, err = os.Create(filepath.Join(directory, fileName))
			if err != nil {
				tErr := transaction.Rollback()
				if tErr != nil {
					return errors.Wrap(tErr, "")
				}
				return fiber.NewError(fiber.StatusInternalServerError, "Error while forming report")
			}
		}
	}
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"Название услуги", "Общая сумма выручки за отчётный период"})
	if err != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while forming report")
	}

	for rows.Next() {
		reportRow := ReportRow{}
		err = rows.Scan(&reportRow.Sum, &reportRow.Name)
		if err != nil {
			log.Println(err)
			continue
		}
		err = writer.Write(reportRow.toString())
		if err != nil {
			log.Println(err)
			continue
		}
	}
	writer.Flush()
	if writer.Error() != nil {
		tErr := transaction.Rollback()
		if tErr != nil {
			return errors.Wrap(tErr, "")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "Error while forming report")
	}
	err = transaction.Commit()
	if err != nil {
		return errors.Wrap(err, "")
	}

	host := c.Hostname()
	url := fmt.Sprintf("%s/report/%s", host, fileName)
	return c.JSON(app.ResponseBody{Data: ReportUrl{Url: url}})
}

func DownloadReport(c *fiber.Ctx) (err error) {
	directory := "./reports"
	fileName := c.Params("filename")
	if fileName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Filename is required")
	}
	_, err = os.Open(filepath.Join(directory, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "Report not found")
		}
		return errors.Wrap(err, "")
	}
	return c.Download(filepath.Join(directory, fileName))
}

func GetOperations(c *fiber.Ctx) (err error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "User Id is required")
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}
	perPage, err := strconv.Atoi(c.Query("per-page"))
	if err != nil {
		perPage = 10
	}
	sort := c.Query("sort")
	if sort == "" {
		sort = "created_at"
	}
	sortDir := c.Query("sort-dir")
	if sortDir == "" {
		sortDir = "asc"
	}
	operationsFilter := OperationsFilter{
		Page:    page,
		PerPage: perPage,
		Sort:    sort,
		SortDir: sortDir,
	}
	errs := validate(operationsFilter)
	if errs != nil {
		return c.
			Status(fiber.StatusUnprocessableEntity).
			JSON(app.ResponseBody{Error: app.Error{Code: fiber.StatusUnprocessableEntity, Fields: errs}})
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
	operation := models.Operation{UserId: user.Id}
	operations, err := operation.GetByUserId(
		transaction,
		operationsFilter.Page,
		operationsFilter.PerPage,
		operationsFilter.Sort,
		operationsFilter.SortDir,
	)
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
	return c.JSON(app.ResponseBody{Data: operations})
}
