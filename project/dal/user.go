package dao

import (
	"GOLANG/project/models"
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type UserDao struct {
	db *sql.DB
}

// InitializeDB initializes the database and creates the users table if it doesn't exist.
func InitializeUserDB() (*UserDao, error) {
	var err error
	dbInstance, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`
	_, err = dbInstance.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return &UserDao{db: dbInstance}, nil
}

// SignUp handles user registration.
func (s *UserDao) AddUser(ctx context.Context, user models.User) error {
	_, err := s.db.Exec("INSERT INTO users (email, password) VALUES (?, ?)", user.Email, user.Pwd)
	if err != nil {
		return err
	}

	return nil
}

// SignIn handles user sign-in.
func (s *UserDao) GetPassword(ctx context.Context, user models.User) (*string, error) {
	var storedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedPassword)
	if err != nil {
		return nil, err
	}

	return &storedPassword, nil
}
