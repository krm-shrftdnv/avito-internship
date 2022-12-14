package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Operation struct {
	Id         int32     `json:"id"`
	UserId     int32     `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	IsIncrease bool      `json:"is_increase"`
	Value      float32   `json:"value"`
}

func (operation *Operation) GetById(transaction *sql.Tx) (*Operation, error) {
	row := transaction.QueryRow("select * from operation where id = $1", operation.Id)
	err := row.Scan(&operation.Id, &operation.UserId, &operation.CreatedAt, &operation.IsIncrease, &operation.Value)
	if err != nil {
		return nil, err
	}
	return operation, nil
}

func (operation *Operation) GetByUserId(transaction *sql.Tx, page int, perPage int, sort string, sortDir string) (operations []Operation, err error) {
	switch sortDir {
	case "desc", "DESC":
		sortDir = "desc"
		break
	default:
		sortDir = "asc"
	}
	switch sort {
	case "created_at":
		sort = "created_at"
		break
	case "value":
		sort = "value"
		break
	default:
		sort = "id"
	}
	var query string
	query = fmt.Sprintf("select * from operation where user_id = $1 order by %s %s limit $2 offset $3", sort, sortDir)
	rows, err := transaction.Query(
		query,
		operation.UserId,
		perPage,
		perPage*(page-1),
	)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		operationModel := Operation{}
		err := rows.Scan(&operationModel.Id, &operationModel.UserId, &operationModel.CreatedAt, &operationModel.IsIncrease, &operationModel.Value)
		if err != nil {
			log.Println(err)
			continue
		}
		operations = append(operations, operationModel)
	}
	return operations, nil
}

func (operation *Operation) Save(transaction *sql.Tx) (err error) {
	if operation.Id == 0 {
		operation.CreatedAt = time.Now()
		err = transaction.
			QueryRow(
				"insert into operation (user_id, created_at, is_increase, value) values ($1, $2, $3, $4) returning id",
				operation.UserId,
				operation.CreatedAt,
				operation.IsIncrease,
				operation.Value,
			).
			Scan(&operation.Id)
		if err != nil {
			return err
		}
	} else {
		operationModel := &Operation{Id: operation.Id}
		operationModel, err = operationModel.GetById(transaction)
		if err == sql.ErrNoRows {
			operation.CreatedAt = time.Now()
			err = transaction.
				QueryRow(
					"insert into operation (id, user_id, created_at, is_increase, value) values ($1, $2, $3, $4, $5) returning id",
					operation.Id,
					operation.UserId,
					operation.CreatedAt,
					operation.IsIncrease,
					operation.Value,
				).
				Scan(&operation.Id)
			if err != nil {
				return err
			}
		} else {
			_, err = transaction.Exec(
				"update operation set user_id = $1, is_increase = $2, value = $3 where id = $4",
				operationModel.UserId,
				operation.IsIncrease,
				operation.Value,
				operation.Id,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
