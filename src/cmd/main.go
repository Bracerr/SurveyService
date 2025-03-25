package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"survey-project/src/internal/config"
	"survey-project/src/internal/database"
	userhttp "survey-project/src/internal/delivery/http"
	"survey-project/src/internal/repository"
	"survey-project/src/internal/usecase"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)

	if err := database.RunMigrations(cfg); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	pool, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(pool)

	userUsecase := usecase.NewUserUsecase(userRepo, refreshTokenRepo, cfg.JWT)

	userHandler := userhttp.NewUserHandler(userUsecase)

	jwtConfig := &userhttp.JWTConfig{Secret: cfg.JWT.Secret}
	router := userhttp.NewRouter(userHandler, logger, jwtConfig)

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Infof("Starting server on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router.Setup()); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
