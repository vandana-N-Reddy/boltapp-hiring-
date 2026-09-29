package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	OTPCode   string    `json:"-"` // never expose in API responses
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CheckoutSubmission struct {
	ID              int       `json:"id"`
	UserID          *int      `json:"userId,omitempty"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	ShippingAddress string    `json:"shippingAddress"`
	CreatedAt       time.Time `json:"createdAt"`
}
