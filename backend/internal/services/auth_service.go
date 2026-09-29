package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"otp-login/internal/models"
	"otp-login/internal/repositories"
)

var ErrInvalidOTP = errors.New("invalid OTP code")

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(email, firstName, lastName string) (*models.User, string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	otp, err := generateOTP()
	if err != nil {
		return nil, "", fmt.Errorf("generate otp: %w", err)
	}

	user, err := s.userRepo.Create(email, firstName, lastName, otp)
	if err != nil {
		return nil, "", err
	}
	return user, otp, nil
}

func (s *AuthService) Recognize(email string) (bool, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	_, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, repositories.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *AuthService) Verify(email, code string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, repositories.ErrNotFound) {
		return nil, ErrInvalidOTP
	}
	if err != nil {
		return nil, err
	}
	if user.OTPCode != code {
		return nil, ErrInvalidOTP
	}
	return user, nil
}

func generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}
