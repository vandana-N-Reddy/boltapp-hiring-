package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"otp-login/internal/repositories"
	"otp-login/internal/services"
)

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type AuthHandler struct {
	authSvc *services.AuthService
}

func NewAuthHandler(authSvc *services.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// POST /api/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if !emailRe.MatchString(req.Email) {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address")
		return
	}
	if req.FirstName == "" {
		writeError(w, http.StatusBadRequest, "First name is required")
		return
	}
	if req.LastName == "" {
		writeError(w, http.StatusBadRequest, "Last name is required")
		return
	}

	user, otp, err := h.authSvc.Register(req.Email, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, repositories.ErrEmailExists) {
			writeError(w, http.StatusConflict, "Email is already registered")
			return
		}
		writeError(w, http.StatusInternalServerError, "Registration failed. Please try again.")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":   "Registration successful",
		"code":      otp,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
}

// GET /api/users/recognize?email=...
func (h *AuthHandler) Recognize(w http.ResponseWriter, r *http.Request) {
	email := strings.TrimSpace(r.URL.Query().Get("email"))
	if email == "" || !emailRe.MatchString(email) {
		writeJSON(w, http.StatusOK, map[string]bool{"registered": false})
		return
	}

	registered, err := h.authSvc.Recognize(email)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Recognition check failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"registered": registered})
}

// POST /api/auth/verify
func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Code = strings.TrimSpace(req.Code)

	if req.Email == "" || !emailRe.MatchString(req.Email) {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address")
		return
	}
	if len(req.Code) != 6 {
		writeError(w, http.StatusBadRequest, "Please enter a 6-digit code")
		return
	}
	for _, c := range req.Code {
		if c < '0' || c > '9' {
			writeError(w, http.StatusBadRequest, "Code must contain digits only")
			return
		}
	}

	user, err := h.authSvc.Verify(req.Email, req.Code)
	if err != nil {
		if errors.Is(err, services.ErrInvalidOTP) {
			writeError(w, http.StatusUnauthorized, "Invalid code. Please try again.")
			return
		}
		writeError(w, http.StatusInternalServerError, "Verification failed. Please try again.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"message":   "Login successful",
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
}
