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

func (order Order) Save() (err error) {

	if order.Id == 0 {

	} else {

	}
	return nil
}

func (order Order) GetById(transaction *sql.Tx) (*Order, error) {
	return nil, nil
}
