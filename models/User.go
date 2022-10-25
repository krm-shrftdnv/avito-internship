package models

import "database/sql"

type User struct {
	Id           int32   `json:"id"`
	Name         string  `json:"name"`
	BalanceValue float32 `json:"balance_value"`
}

func (user User) GetById(transaction *sql.Tx) (*User, error) {
	row := transaction.QueryRow("select * from \"user\" where id = $1", user.Id)
	err := row.Scan(&user.Id, &user.Name, &user.BalanceValue)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
