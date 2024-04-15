package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Users struct {
	db *sql.DB
}

type UserEntry struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateTable() *Users {
	const create string = `
  CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  phone VARCHAR(50) NOT NULL
  );`
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatalf("Failed to open users.db:\n%s", err)
	}
	if _, err := db.Exec(create); err != nil {
		log.Fatalf("Failed to create table in database:\n%s", err)
	}
	log.Println("Successfully loaded users.db")
	return &Users{
		db: db,
	}
}

func (u *Users) ListAll() ([]*UserEntry, error) {
	var users []*UserEntry
	const listAll string = `
  SELECT * FROM users`
	rows, err := u.db.Query(listAll)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var entry UserEntry
		if err := rows.Scan(&entry.Name, &entry.Phone); err != nil {
			return nil, err
		}
		users = append(users, &entry)
	}
	return users, nil
}

func (u *Users) Append(entry *UserEntry) error {
	const insert string = `
  INSERT INTO users VALUES(?, ?)`
	if _, err := u.db.Exec(insert, entry.Name, entry.Phone); err != nil {
		return err
	}
	return nil
}

func (u *Users) DeleteByName(name string) (bool, error) {
	const deleteByName string = `
  DELETE FROM users
  WHERE name=?`
	res, err := u.db.Exec(deleteByName, name)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (u *Users) DeleteByPhone(phone string) (bool, error) {
	const deleteByPhone string = `
  DELETE FROM users
  WHERE phone=?`
	res, err := u.db.Exec(deleteByPhone, phone)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
