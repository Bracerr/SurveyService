package http

import (
	"survey-project/src/internal/config"
	"survey-project/src/pkg/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

type Router struct {
	userHandler *UserHandler
	logger      *logrus.Logger
	jwtConfig   *config.JWTConfig
}

func NewRouter(userHandler *UserHandler, logger *logrus.Logger, jwtConfig *config.JWTConfig) *Router {
	return &Router{
		userHandler: userHandler,
		logger:      logger,
		jwtConfig:   jwtConfig,
	}
}

func (r *Router) Setup() *chi.Mux {
	router := chi.NewRouter()

	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	router.Route("/api/v1", func(router chi.Router) {
		router.Route("/auth", func(router chi.Router) {
			router.Post("/register", r.userHandler.Register)
			router.Post("/login", r.userHandler.Login)
			router.Post("/refresh", r.userHandler.RefreshToken)
		})

		router.Route("/users", func(router chi.Router) {
			router.Use(middleware.AuthMiddleware(r.jwtConfig.Secret))

			router.Get("/", r.userHandler.GetAll)
			router.Get("/{id}", r.userHandler.GetByID)
			router.Put("/{id}", r.userHandler.Update)
			router.Delete("/{id}", r.userHandler.Delete)
		})
	})

	return router
}
