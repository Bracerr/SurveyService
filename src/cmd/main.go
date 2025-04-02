package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"survey-project/src/internal/bootstrap"
	"survey-project/src/internal/config"
	userhttp "survey-project/src/internal/delivery/http"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := bootstrap.InitLogger()

	postgresDB, mongoDB, err := bootstrap.InitDatabases(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize databases: %v", err)
	}
	defer postgresDB.Close()
	defer mongoDB.Disconnect(context.Background())

	userRepo, refreshTokenRepo, surveyRepo := bootstrap.InitRepositories(postgresDB, mongoDB)

	userUsecase, surveyUsecase := bootstrap.InitUseCases(userRepo, refreshTokenRepo, surveyRepo, cfg.JWT)

	userHandler, surveyHandler := bootstrap.InitHandlers(userUsecase, surveyUsecase)

	router := userhttp.NewRouter(userHandler, surveyHandler, logger, &cfg.JWT)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router.SetupRoutes(),
	}

	bootstrap.SetupGracefulShutdown(server, logger)

	logger.Infof("Starting server on %s", fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}
