package http

import (
	"net/http"

	"survey-project/src/internal/config"
	"survey-project/src/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Router struct {
	userHandler   *UserHandler
	surveyHandler *SurveyHandler
	logger        *logrus.Logger
	jwtConfig     *config.JWTConfig
}

func NewRouter(userHandler *UserHandler, surveyHandler *SurveyHandler, logger *logrus.Logger, jwtConfig *config.JWTConfig) *Router {
	return &Router{
		userHandler:   userHandler,
		surveyHandler: surveyHandler,
		logger:        logger,
		jwtConfig:     jwtConfig,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.Logger)

	// Public routes
	router.Group(func(router chi.Router) {
		router.Post("/api/v1/auth/register", r.userHandler.Register)
		router.Post("/api/v1/auth/login", r.userHandler.Login)
		router.Post("/api/v1/auth/refresh", r.userHandler.RefreshToken)
	})

	// Protected routes
	router.Group(func(router chi.Router) {
		router.Use(middleware.AuthMiddleware(r.jwtConfig.Secret))

		// User routes
		router.Route("/api/v1/users", func(router chi.Router) {
			// Routes for current user
			router.Get("/me", r.userHandler.GetMe)
			router.Put("/me", r.userHandler.UpdateMe)
			router.Delete("/me", r.userHandler.DeleteMe)

			// Admin routes
			router.Group(func(router chi.Router) {
				router.Use(middleware.AdminOnly)
				router.Get("/", r.userHandler.GetAll)
				router.Get("/{id}", r.userHandler.GetByID)
				router.Put("/{id}", r.userHandler.Update)
				router.Delete("/{id}", r.userHandler.Delete)
			})
		})

		// Survey routes
		router.Route("/api/v1/surveys", func(router chi.Router) {
			router.Group(func(router chi.Router) {
				router.Use(middleware.AdminOnly)
				router.Get("/", r.surveyHandler.GetAll)
				router.Get("/user/{user_id}", r.surveyHandler.GetByUserID)
			})

			router.Post("/", r.surveyHandler.Create)
			router.Get("/{id}", r.surveyHandler.GetByID)
			router.Put("/{id}", r.surveyHandler.Update)
			router.Delete("/{id}", r.surveyHandler.Delete)
			router.Get("/my", r.surveyHandler.GetMy)
		})
	})

	return router
}
