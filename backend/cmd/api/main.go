package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	appAuth "legalflow/internal/application/auth"
	appClient "legalflow/internal/application/client"
	"legalflow/internal/application/user"
	"legalflow/internal/config"
	"legalflow/internal/database"
	"legalflow/internal/health"
	httpAuth "legalflow/internal/http/auth"
	httpClient "legalflow/internal/http/client"
	httpuser "legalflow/internal/http/user"
	"legalflow/internal/infrastructure/auth"
	"legalflow/internal/infrastructure/hasher"
	"legalflow/internal/infrastructure/persistence/postgres"
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

	if db != nil {
		userRepo := postgres.NewUserRepository(db)
		pwdHasher := hasher.NewBcryptHasher()
		pwdVerifier := hasher.NewBcryptVerifier()
		jwtService := auth.NewJWTService(cfg.JWTSecret, cfg.JWTExpirationMinutes)

		createUser := user.NewCreateUserUseCase(userRepo, pwdHasher)
		userHandler := httpuser.NewHandler(createUser)
		userHandler.Register(r)

		authUseCase := appAuth.NewAuthenticateUserUseCase(userRepo, pwdVerifier)
		loginHandler := httpAuth.NewLoginHandler(authUseCase, jwtService)
		loginHandler.Register(r)

		getCurrentUser := appAuth.NewGetCurrentUserUseCase(userRepo)
		meHandler := httpAuth.NewMeHandler(getCurrentUser)

		clientRepo := postgres.NewClientRepository(db)
		createClient := appClient.NewCreateClientUseCase(clientRepo)
		clientHandler := httpClient.NewHandler(createClient)

		listClients := appClient.NewListClientsUseCase(clientRepo)
		listHandler := httpClient.NewListHandler(listClients)

		getClient := appClient.NewGetClientByIDUseCase(clientRepo)
		getHandler := httpClient.NewGetHandler(getClient)

		updateClient := appClient.NewUpdateClientUseCase(clientRepo)
		updateHandler := httpClient.NewUpdateHandler(updateClient)

		authMW := middleware.AuthMiddleware(newJWTAdapter(jwtService))
		r.Group(func(r chi.Router) {
			r.Use(authMW)
			meHandler.Register(r)
			clientHandler.Register(r)
			listHandler.Register(r)
			getHandler.Register(r)
			updateHandler.Register(r)
		})
	}

	log.Info("Server starting on :" + cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, r); err != nil {
		panic(err)
	}
}

type jwtAdapter struct {
	svc *auth.JWTService
}

func (a *jwtAdapter) Validate(token string) (*middleware.TokenClaims, error) {
	claims, err := a.svc.Validate(token)
	if err != nil {
		return nil, err
	}
	return &middleware.TokenClaims{
		UserID:    claims.UserID,
		Email:     claims.Email,
		IssuedAt:  claims.IssuedAt,
		ExpiresAt: claims.ExpiresAt,
	}, nil
}

func newJWTAdapter(svc *auth.JWTService) *jwtAdapter {
	return &jwtAdapter{svc: svc}
}
