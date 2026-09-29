package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"otp-login/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// We test the handler layer directly using httptest.
// Validation fires before any DB/service call, so nil services are safe for these tests.

// Helper to post JSON and return response
func postJSON(t *testing.T, router http.Handler, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func getRequest(t *testing.T, router http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// --- Validation-only tests (no DB required) ---

func TestRegister_EmptyEmail(t *testing.T) {
	// Use nil services — validation fires before DB call
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/register", authH.Register)

	rr := postJSON(t, r, "/api/auth/register", map[string]string{
		"email": "", "firstName": "John", "lastName": "Doe",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/register", authH.Register)

	rr := postJSON(t, r, "/api/auth/register", map[string]string{
		"email": "not-an-email", "firstName": "John", "lastName": "Doe",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRegister_EmptyFirstName(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/register", authH.Register)

	rr := postJSON(t, r, "/api/auth/register", map[string]string{
		"email": "john@example.com", "firstName": "", "lastName": "Doe",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRegister_EmptyLastName(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/register", authH.Register)

	rr := postJSON(t, r, "/api/auth/register", map[string]string{
		"email": "john@example.com", "firstName": "John", "lastName": "",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestVerify_ShortCode(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/verify", authH.Verify)

	rr := postJSON(t, r, "/api/auth/verify", map[string]string{
		"email": "john@example.com", "code": "123",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestVerify_NonNumericCode(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/auth/verify", authH.Verify)

	rr := postJSON(t, r, "/api/auth/verify", map[string]string{
		"email": "john@example.com", "code": "12345a",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCheckout_EmptyEmail(t *testing.T) {
	checkoutH := handlers.NewCheckoutHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/checkout", checkoutH.Submit)

	rr := postJSON(t, r, "/api/checkout", map[string]string{
		"email": "", "phone": "9876543210", "shippingAddress": "Bengaluru",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCheckout_InvalidPhone(t *testing.T) {
	checkoutH := handlers.NewCheckoutHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/checkout", checkoutH.Submit)

	rr := postJSON(t, r, "/api/checkout", map[string]string{
		"email": "john@example.com", "phone": "abc", "shippingAddress": "Bengaluru",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCheckout_EmptyAddress(t *testing.T) {
	checkoutH := handlers.NewCheckoutHandler(nil)
	r := chi.NewRouter()
	r.Post("/api/checkout", checkoutH.Submit)

	rr := postJSON(t, r, "/api/checkout", map[string]string{
		"email": "john@example.com", "phone": "9876543210", "shippingAddress": "",
	})
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestRecognize_InvalidEmail(t *testing.T) {
	authH := handlers.NewAuthHandler(nil)
	r := chi.NewRouter()
	r.Get("/api/users/recognize", authH.Recognize)

	rr := getRequest(t, r, "/api/users/recognize?email=notvalid")
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
	var resp map[string]bool
	json.NewDecoder(rr.Body).Decode(&resp)
	if resp["registered"] {
		t.Error("expected registered=false for invalid email")
	}
}

// Note: Tests requiring a real DB connection are in backend/internal/integration_test.go
// Run with: go test ./... -tags integration
