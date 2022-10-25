package models

import "database/sql"

type User struct {
	Id           int32   `json:"id"`
	Name         string  `json:"name"`
	BalanceValue float32 `json:"balance_value"`
}

func (user User) GetById(transaction *sql.Tx) (*User, error) {
	rows, err := transaction.Query("select * from user where id = ?", user.Id)
	if err != nil {
		return nil, err
	}
	err = rows.Scan(&user.Id, &user.Name, &user.BalanceValue)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
