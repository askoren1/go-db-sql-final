package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуем добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}
	// верните идентификатор последней добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	row := s.db.QueryRow("SELECT * FROM parcel WHERE Number = :Number", sql.Named("Number", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуем чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT *  FROM parcel WHERE Client = :Client", sql.Named("Client", client))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for rows.Next() {
		r := Parcel{}

		err := rows.Scan(&r.Number, &r.Client, &r.Status, &r.Address, &r.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		res = append(res, r)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET Status = :Status WHERE Number = :Number",
		sql.Named("Status", status),
		sql.Named("Number", number))
	return err

}

func (s ParcelStore) SetAddress(number int, address string) error {

	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE Number = :Number AND Status = :Status",
		sql.Named("address", address),
		sql.Named("Number", number),
		sql.Named("Status", "registered"))
	return err

}

func (s ParcelStore) Delete(number int) error {
	// Удалим строку из таблицы parcel
	_, err := s.db.Exec("DELETE FROM parcel WHERE Number = :Number AND Status = :Status",
		sql.Named("Number", number),
		sql.Named("Status", "registered"))
	if err != nil {
		return err
	}

	return nil
}
