package models

import (
	"database/sql"
	"time"
)

type Reserve struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id"`
	OrderId   int32     `json:"order_id"`
	Value     float32   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	IsDebited bool      `json:"is_debited"`
}

func (reserve Reserve) GetById(transaction *sql.Tx) (*Reserve, error) {
	row := transaction.QueryRow("select * from reserve where id = $1", reserve.Id)
	err := row.Scan(&reserve.Id, &reserve.UserId, &reserve.OrderId, &reserve.Value, &reserve.CreatedAt, &reserve.IsDebited)
	if err != nil {
		return nil, err
	}
	return &reserve, nil
}

func (reserve Reserve) Save(transaction *sql.Tx) (err error) {
	if reserve.Id == 0 {
		reserve.IsDebited = false
		reserve.CreatedAt = time.Now()
		err = transaction.QueryRow(
			"insert into reserve (user_id, order_id, value, created_at, is_debited) values ($1, $2, $3, $4, $5) returning id",
			reserve.UserId,
			reserve.OrderId,
			reserve.Value,
			reserve.CreatedAt,
			reserve.IsDebited,
		).Scan(&reserve.Id)
		if err != nil {
			return err
		}
	} else {
		reserveModel, err := reserve.GetById(transaction)
		if err == sql.ErrNoRows {
			reserve.IsDebited = false
			reserve.CreatedAt = time.Now()
			err = transaction.QueryRow(
				"insert into reserve (id, user_id, order_id, value, created_at, is_debited) values ($1, $2, $3, $4, $5, $6) returning id",
				reserve.Id,
				reserve.UserId,
				reserve.OrderId,
				reserve.Value,
				reserve.CreatedAt,
				reserve.IsDebited,
			).Scan(&reserve.Id)
			if err != nil {
				return err
			}
		} else {
			_, err = transaction.Exec(
				"update reserve set user_id = $1, order_id = $2, value=$3, is_debited=$4 where id = $5",
				reserve.UserId,
				reserve.OrderId,
				reserve.Value,
				reserve.IsDebited,
				reserveModel.Id,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
