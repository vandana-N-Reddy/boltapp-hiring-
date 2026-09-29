package services

import (
	"errors"
	"strings"

	"otp-login/internal/models"
	"otp-login/internal/repositories"
)

type CheckoutService struct {
	checkoutRepo *repositories.CheckoutRepository
	userRepo     *repositories.UserRepository
}

func NewCheckoutService(checkoutRepo *repositories.CheckoutRepository, userRepo *repositories.UserRepository) *CheckoutService {
	return &CheckoutService{checkoutRepo: checkoutRepo, userRepo: userRepo}
}

func (s *CheckoutService) Submit(email, phone, shippingAddress string) (*models.CheckoutSubmission, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	phone = strings.TrimSpace(phone)
	shippingAddress = strings.TrimSpace(shippingAddress)

	var userID *int
	user, err := s.userRepo.FindByEmail(email)
	if err == nil {
		userID = &user.ID
	} else if !errors.Is(err, repositories.ErrNotFound) {
		return nil, err
	}

	return s.checkoutRepo.Create(email, phone, shippingAddress, userID)
}
