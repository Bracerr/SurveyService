package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"survey-project/src/internal/config"
	"survey-project/src/internal/database"
	userhttp "survey-project/src/internal/delivery/http"
	"survey-project/src/internal/repository"
	"survey-project/src/internal/usecase"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}

func InitDatabases(
	cfg *config.Config,
) (*database.PostgresDB, *database.MongoDB, error) {
	if err := database.RunMigrations(cfg); err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	postgresDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	mongoDB, err := database.NewMongoDB(cfg.MongoDB.URI)
	if err != nil {
		postgresDB.Close()
		return nil, nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	return postgresDB, mongoDB, nil
}

func InitRepositories(
	postgresDB *database.PostgresDB,
	mongoDB *database.MongoDB,
) (*repository.UserRepository, *repository.RefreshTokenRepository, *repository.SurveyRepository) {
	userRepo := repository.NewUserRepository(postgresDB.Pool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(postgresDB.Pool)
	surveyRepo := repository.NewSurveyRepository(
		mongoDB.Database("survey").Collection("surveys"),
	)
	return userRepo, refreshTokenRepo, surveyRepo
}

func InitUseCases(
	userRepo *repository.UserRepository,
	refreshTokenRepo *repository.RefreshTokenRepository,
	surveyRepo *repository.SurveyRepository,
	jwtConfig config.JWTConfig,
) (*usecase.UserUsecase, *usecase.SurveyUsecase) {
	userUsecase := usecase.NewUserUsecase(*userRepo, *refreshTokenRepo, jwtConfig)
	surveyUsecase := usecase.NewSurveyUsecase(surveyRepo)
	return userUsecase, surveyUsecase
}

func InitHandlers(
	userUsecase *usecase.UserUsecase,
	surveyUsecase *usecase.SurveyUsecase,
) (*userhttp.UserHandler, *userhttp.SurveyHandler) {
	userHandler := userhttp.NewUserHandler(userUsecase)
	surveyHandler := userhttp.NewSurveyHandler(surveyUsecase)
	return userHandler, surveyHandler
}

func SetupGracefulShutdown(server *http.Server, logger *logrus.Logger) {
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Errorf("Server forced to shutdown: %v", err)
		}
	}()
}
