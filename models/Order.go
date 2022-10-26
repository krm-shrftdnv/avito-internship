package models

import (
	"database/sql"
	"time"
)

type Order struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id"`
	ServiceId int32     `json:"service_id"`
	CreatedAt time.Time `json:"created_at"`
	IsPaid    bool      `json:"is_paid"`
}

func (order Order) GetById(transaction *sql.Tx) (*Order, error) {
	row := transaction.QueryRow("select * from \"order\" where id = $1", order.Id)
	err := row.Scan(&order.Id, &order.UserId, &order.ServiceId, &order.CreatedAt, &order.IsPaid)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (order Order) Save(transaction *sql.Tx) (err error) {
	if order.Id == 0 {
		order.IsPaid = false
		order.CreatedAt = time.Now()
		var id int
		err = transaction.QueryRow(
			"insert into \"order\" (user_id, service_id, created_at, is_paid) values ($1, $2, $3, $4) returning id",
			order.UserId,
			order.ServiceId,
			order.CreatedAt,
			order.IsPaid,
		).Scan(&id)
		if err != nil {
			return err
		}
		order.Id = int32(id)
	} else {
		orderModel, err := order.GetById(transaction)
		if err == sql.ErrNoRows {
			order.IsPaid = false
			order.CreatedAt = time.Now()
			err = transaction.QueryRow(
				"insert into \"order\" (id, user_id, service_id, created_at, is_paid) values ($1, $2, $3, $4, $5) returning id",
				order.Id,
				order.UserId,
				order.ServiceId,
				order.CreatedAt,
				order.IsPaid,
			).Scan(&order.Id)
			if err != nil {
				return err
			}
		} else {
			_, err = transaction.Exec(
				"update \"order\" set user_id = $1, service_id = $2, is_paid=$3 where id = $4",
				order.UserId,
				order.ServiceId,
				order.IsPaid,
				orderModel.Id,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
