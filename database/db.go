package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type PhoneBook struct {
	db *sql.DB
}

type UserEntry struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateTable(path string) *PhoneBook {
	const create string = `
  CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(120) NOT NULL,
  phone VARCHAR(50) NOT NULL
  );`
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Fatalf("Failed to open users.db:\n%s", err)
	}
	if _, err := db.Exec(create); err != nil {
		log.Fatalf("Failed to create table in database:\n%s", err)
	}
	log.Println("Successfully loaded users.db")
	return &PhoneBook{
		db: db,
	}
}

func (p *PhoneBook) ListAll() ([]*UserEntry, error) {
	var users []*UserEntry
	const listAll string = `
  SELECT * FROM users;`
	rows, err := p.db.Query(listAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var entry UserEntry
		if err := rows.Scan(&id, &entry.Name, &entry.Phone); err != nil {
			return nil, err
		}
		users = append(users, &entry)
	}
	return users, nil
}

func (p *PhoneBook) Append(entry *UserEntry) error {
	const insert string = `
  INSERT INTO users(name, phone) VALUES (?, ?);`
	if _, err := p.db.Exec(insert, entry.Name, entry.Phone); err != nil {
		return err
	}
	return nil
}

func (p *PhoneBook) DeleteByName(name string) (bool, error) {
	const deleteByName string = `
  DELETE FROM users
  WHERE name=?;`
	res, err := p.db.Exec(deleteByName, name)
	if err != nil {
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}

func (p *PhoneBook) DeleteByPhone(phone string) (string, bool, error) {
	const find string = `SELECT name FROM users WHERE phone=?;`
	var name string
	if err := p.db.QueryRow(find, phone).Scan(&name); err != nil {
		return "", false, err
	}
	const deleteByPhone string = `
  DELETE FROM users
  WHERE phone=?;`
	res, err := p.db.Exec(deleteByPhone, phone)
	if err != nil {
		return "", false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "", false, err
	}
	return name, rows > 0, nil
}

func (p *PhoneBook) Close() {
	p.db.Close()
}
