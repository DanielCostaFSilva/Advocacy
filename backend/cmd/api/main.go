package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"legalflow/internal/database"
	"legalflow/internal/health"
)

func main() {
	cfg := database.ConfigFromEnv()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Printf("Warning: database connection failed: %v", err)
	} else {
		defer db.Close()
		log.Println("Database connected successfully")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	healthHandler := health.NewHandler()
	healthHandler.Register(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
