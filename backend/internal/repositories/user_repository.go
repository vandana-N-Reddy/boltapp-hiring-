package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"otp-login/internal/models"
)

var ErrEmailExists = errors.New("email already registered")
var ErrNotFound = errors.New("user not found")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(email, firstName, lastName, otpCode string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`INSERT INTO users (email, first_name, last_name, otp_code)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, email, first_name, last_name, created_at, updated_at`,
		email, firstName, lastName, otpCode,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrEmailExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		`SELECT id, email, first_name, last_name, otp_code, created_at, updated_at
		 FROM users WHERE email = $1`, email,
	).Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.OTPCode, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("find user: %w", err)
	}
	return &u, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
