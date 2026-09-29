package repositories

import (
	"database/sql"
	"fmt"

	"otp-login/internal/models"
)

type CheckoutRepository struct {
	db *sql.DB
}

func NewCheckoutRepository(db *sql.DB) *CheckoutRepository {
	return &CheckoutRepository{db: db}
}

func (r *CheckoutRepository) Create(email, phone, shippingAddress string, userID *int) (*models.CheckoutSubmission, error) {
	var c models.CheckoutSubmission
	err := r.db.QueryRow(
		`INSERT INTO checkout_submissions (user_id, email, phone, shipping_address)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, email, phone, shipping_address, created_at`,
		userID, email, phone, shippingAddress,
	).Scan(&c.ID, &c.UserID, &c.Email, &c.Phone, &c.ShippingAddress, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create checkout: %w", err)
	}
	return &c, nil
}
