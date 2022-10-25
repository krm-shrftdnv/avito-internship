package models

type Service struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}
