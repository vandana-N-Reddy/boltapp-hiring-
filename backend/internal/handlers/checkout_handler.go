package handlers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"otp-login/internal/services"
)

var phoneRe = regexp.MustCompile(`^\+?[0-9\s\-().]{7,20}$`)

type CheckoutHandler struct {
	checkoutSvc *services.CheckoutService
}

func NewCheckoutHandler(checkoutSvc *services.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{checkoutSvc: checkoutSvc}
}

// POST /api/checkout
func (h *CheckoutHandler) Submit(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		ShippingAddress string `json:"shippingAddress"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)
	req.ShippingAddress = strings.TrimSpace(req.ShippingAddress)

	if req.Email == "" {
		writeError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if !emailRe.MatchString(req.Email) {
		writeError(w, http.StatusBadRequest, "Please enter a valid email address")
		return
	}
	if req.Phone == "" {
		writeError(w, http.StatusBadRequest, "Phone number is required")
		return
	}
	if !phoneRe.MatchString(req.Phone) {
		writeError(w, http.StatusBadRequest, "Please enter a valid phone number")
		return
	}
	if req.ShippingAddress == "" {
		writeError(w, http.StatusBadRequest, "Shipping address is required")
		return
	}

	_, err := h.checkoutSvc.Submit(req.Email, req.Phone, req.ShippingAddress)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Checkout submission failed. Please try again.")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": "Checkout submitted successfully",
	})
}
