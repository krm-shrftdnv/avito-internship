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
	//row := transaction.QueryRow("select * from \"order\" where id = $1", order.Id)
	//err := row.Scan(&order.Id, &order.UserId, &order.OrderId, &order.Value, &order.CreatedAt, &order.IsDebited)
	//if err != nil {
	//	return nil, err
	//}
	//return &order, nil
	return nil, nil
}

func (order Order) Save(transaction *sql.Tx) (err error) {
	//if order.Id == 0 {
	//	order.IsDebited = false
	//	order.CreatedAt = time.Now()
	//	_, err := transaction.Exec(
	//		"insert into \"order\" (user_id, order_id, value, created_at, is_debited) values ($1, $2, $3, $4, $5)",
	//		order.UserId,
	//		order.OrderId,
	//		order.Value,
	//		order.CreatedAt,
	//		order.IsDebited,
	//	)
	//	if err != nil {
	//		return err
	//	}
	//} else {
	//	reserveModel, err := order.GetById(transaction)
	//	if err == sql.ErrNoRows {
	//		order.IsDebited = false
	//		order.CreatedAt = time.Now()
	//		_, err := transaction.Exec(
	//			"insert into \"order\" (id, user_id, order_id, value, created_at, is_debited) values ($1, $2, $3, $4, $5, $6)",
	//			order.Id,
	//			order.UserId,
	//			order.OrderId,
	//			order.Value,
	//			order.CreatedAt,
	//			order.IsDebited,
	//		)
	//		if err != nil {
	//			return err
	//		}
	//	} else {
	//		_, err = transaction.Exec(
	//			"update \"order\" set user_id = $1, order_id = $2, value=$3, is_debited=$4 where id = $5",
	//			order.UserId,
	//			order.OrderId,
	//			order.Value,
	//			order.IsDebited,
	//			reserveModel.Id,
	//		)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}
	return nil
}
