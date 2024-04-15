package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type PhoneBook struct {
	db *sql.DB
}

type UserEntry struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func CreateTable(path string) (*PhoneBook, error) {
	const create string = `
  CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(256) NOT NULL,
  phone VARCHAR(64) NOT NULL
  );`
	os.MkdirAll(filepath.Dir(path), os.ModePerm)
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(create); err != nil {
		return nil, err
	}
	return &PhoneBook{
		db: db,
	}, nil
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
		if err == sql.ErrNoRows {
			return "", false, nil
		}
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
