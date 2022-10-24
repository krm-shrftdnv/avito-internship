package models

import "time"

type Reserve struct {
	Id        int32
	UserId    int32
	OrderId   int32
	Value     float32
	CreatedAt time.Time
	IsDebited bool
}
