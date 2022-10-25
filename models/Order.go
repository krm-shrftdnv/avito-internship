package models

import "time"

type Order struct {
	Id        int32     `json:"id"`
	UserId    int32     `json:"user_id"`
	ServiceId int32     `json:"service_id"`
	CreatedAt time.Time `json:"created_at"`
	IsPaid    bool      `json:"is_paid"`
}
