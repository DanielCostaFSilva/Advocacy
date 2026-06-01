package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"legalflow/internal/config"
	"legalflow/internal/database"
	"legalflow/internal/health"
	"legalflow/internal/logger"
	"legalflow/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.NewDefault()

	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	})
	if err != nil {
		log.Info("Warning: database connection failed: " + err.Error())
	} else {
		defer db.Close()
		log.Info("Database connected successfully")
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recovery(log.Inner()))
	r.Use(middleware.AccessLog(log.Inner()))

	healthHandler := health.NewHandler()
	healthHandler.Register(r)

	log.Info("Server starting on :" + cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		panic(err)
	}
}
