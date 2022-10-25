package models

import "time"

type Reserve struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id"`
	OrderId   int32     `json:"order_id"`
	Value     float32   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	IsDebited bool      `json:"is_debited"`
}
