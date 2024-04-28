package cmd

import (
	"context"
	"fitnessme/cors"
	"fitnessme/usermanagement/pkg/db"
	"fitnessme/usermanagement/pkg/repository"
	"fitnessme/usermanagement/pkg/services"

	"fmt"
	"net/http"

	"fitnessme/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type UserApp struct {
	router http.Handler
	db.Handler
	repo repository.UserRepository
	jw   utils.JwtWrapper
}

func New(jwtWrapper *utils.JwtWrapper) *UserApp {
	handler := db.Init("localhost", "users", "postgres", "123", 5432)
	repo := repository.NewUserRepository(handler)

	app := &UserApp{
		router: loadRoutes(repo, jwtWrapper),
		repo:   repo,
		jw:     *jwtWrapper,
	}

	app.router = cors.Cors(app.router)
	return app
}

func loadRoutes(repo repository.UserRepository, jwtWrapper *utils.JwtWrapper) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/user", func(router chi.Router) {
		loadAuthRoutes(router, repo, jwtWrapper)
	})

	return router
}

func loadAuthRoutes(router chi.Router, repo repository.UserRepository, jwtWrapper *utils.JwtWrapper) {
	userManagementService := services.NewUserManagementService(repo, jwtWrapper)
	router.Post("/register", userManagementService.Register)
	router.Post("/", userManagementService.Login)
	router.Get("/", userManagementService.GetUser)
	router.Patch("/", userManagementService.UpdateUser)
	router.Get("/admins", userManagementService.GetAllAdmins)
	router.Get("/name", userManagementService.GetNameById)
	router.Get("/emails", userManagementService.GetAllEmails)
}

func (a *UserApp) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3000",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
