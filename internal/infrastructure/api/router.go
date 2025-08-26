package api

import (
	"net/http"

	"github.com/base-intern-august-b/clipboard-server/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	userUsecase usecase.UserUsecase
}

func NewRouter(userUsecase usecase.UserUsecase) *Router {
	return &Router{
		userUsecase: userUsecase,
	}
}

func (r *Router) Setup() http.Handler {
	router := chi.NewRouter()

	// ミドルウェアの設定
	router.Use(LoggingMiddleware)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Route("/api/v1", func(v1 chi.Router) {
		// ユーザーAPI
		userHandler := NewUserHandler(r.userUsecase)
		v1.Route("/users", func(user chi.Router) {
			user.Post("/", userHandler.CreateUser)
			user.Get("/", userHandler.GetUsers)
			user.Get("/{userID}", userHandler.GetUserByID)
			user.Patch("/{userID}", userHandler.PatchUser)
			user.Post("/{userID}/change-password", userHandler.ChangePassword)
			user.Delete("/{userID}", userHandler.DeleteUser)
		})
	})

	return router
}
