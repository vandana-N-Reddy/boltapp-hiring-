package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"otp-login/internal/database"
	"otp-login/internal/handlers"
	"otp-login/internal/repositories"
	"otp-login/internal/services"
)

func main() {
	// Load .env in development; ignore error in production
	_ = godotenv.Load()

	db, err := database.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// Wire up layers
	userRepo := repositories.NewUserRepository(db)
	checkoutRepo := repositories.NewCheckoutRepository(db)

	authSvc := services.NewAuthService(userRepo)
	checkoutSvc := services.NewCheckoutService(checkoutRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authSvc)
	checkoutHandler := handlers.NewCheckoutHandler(checkoutSvc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   getAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/verify", authHandler.Verify)
		r.Get("/users/recognize", authHandler.Recognize)
		r.Post("/checkout", checkoutHandler.Submit)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getAllowedOrigins() []string {
	origin := os.Getenv("ALLOWED_ORIGIN")
	if origin == "" {
		return []string{"http://localhost:3000", "http://localhost:5500", "http://127.0.0.1:5500"}
	}
	return []string{origin, "http://localhost:3000", "http://localhost:5500"}
}
