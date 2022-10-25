package models

type User struct {
	Id           int32   `json:"id"`
	Name         string  `json:"name"`
	BalanceValue float32 `json:"balance_value"`
}
