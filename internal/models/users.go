package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO users (name, email, hashed_password)
	VALUES($1, $2, $3)
	`
	_, err = m.DB.Exec(context.Background(), query, name, email, string(hashedPassword))
	if err != nil {
		if err, ok := err.(*pgconn.PgError); ok && err.Code == "23505" { // unique violation
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	query := "SELECT id, hashed_password FROM users WHERE email = $1;"

	err := m.DB.QueryRow(context.Background(), query, email).Scan(&id, &hashedPassword)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, ErrInvalidCredentials
		}
		return 0, err

	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
