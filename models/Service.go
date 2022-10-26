package models

import (
	"database/sql"
)

type Service struct {
	Id    int32   `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func (service Service) GetById(transaction *sql.Tx) (*Service, error) {
	row := transaction.QueryRow("select * from service where id = $1", service.Id)
	err := row.Scan(&service.Id, &service.Name, &service.Price)
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (service Service) Save(transaction *sql.Tx) (err error) {
	if service.Id == 0 {
		err = transaction.
			QueryRow("insert into service (name, price) values ($1, $2) returning id", service.Name, service.Price).
			Scan(&service.Id)
		if err != nil {
			return err
		}
	} else {
		userModel, err := service.GetById(transaction)
		if err == sql.ErrNoRows {
			err = transaction.
				QueryRow("insert into service (id, name, price) values ($1, $2, $3) returning id", service.Id, service.Name, service.Price).
				Scan(&service.Id)
			if err != nil {
				return err
			}
		} else {
			_, err = transaction.Exec("update service set name = $1, price = $2 where id = $3", userModel.Name, service.Price, service.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
