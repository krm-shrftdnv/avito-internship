package models

import "time"

type Order struct {
	Id        int32
	UserId    int32
	ServiceId int32
	CreatedAt time.Time
	IsPaid    bool
}
