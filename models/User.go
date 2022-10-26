package models

import "database/sql"

type User struct {
	Id           int32   `json:"id"`
	Name         string  `json:"name"`
	BalanceValue float32 `json:"balance"`
}

func (user User) GetById(transaction *sql.Tx) (*User, error) {
	row := transaction.QueryRow("select * from \"user\" where id = $1", user.Id)
	err := row.Scan(&user.Id, &user.Name, &user.BalanceValue)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (user User) Save(transaction *sql.Tx) (err error) {
	if user.Id == 0 {
		err = transaction.
			QueryRow("insert into \"user\" (name, balance_value) values ($1, $2) returning id", user.Name, user.BalanceValue).
			Scan(&user.Id)
		if err != nil {
			return err
		}
	} else {
		userModel, err := user.GetById(transaction)
		if err == sql.ErrNoRows {
			err = transaction.
				QueryRow("insert into \"user\" (id, name, balance_value) values ($1, $2, $3) returning id", user.Id, user.Name, user.BalanceValue).
				Scan(&user.Id)
			if err != nil {
				return err
			}
		} else {
			_, err = transaction.Exec("update \"user\" set name = $1, balance_value = $2 where id = $3", userModel.Name, user.BalanceValue, user.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
