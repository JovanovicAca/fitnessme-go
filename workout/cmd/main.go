package cmd

import (
	"context"
	"fitnessme/cors"
	"fitnessme/utils"
	"fitnessme/workout/pkg/db"
	"fitnessme/workout/pkg/repository"
	"fitnessme/workout/pkg/services"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WorkoutApp struct {
	router http.Handler
	db.Handler
	repo repository.WorkoutRepository
	jw   utils.JwtWrapper
}

func New(jwtWrapper *utils.JwtWrapper) *WorkoutApp {
	handler := db.Init("localhost", "workout", "postgres", "123", 5432)
	repo := repository.NewWorkoutRepository(handler)

	app := &WorkoutApp{
		router: loadRoutes(repo, jwtWrapper),
		repo:   repo,
		jw:     *jwtWrapper,
	}

	app.router = cors.Cors(app.router)
	return app
}

func loadRoutes(repo repository.WorkoutRepository, jwtWrapper *utils.JwtWrapper) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/workout", func(router chi.Router) {
		loadAuthRoutes(router, repo, jwtWrapper)
	})

	return router
}

func loadAuthRoutes(router chi.Router, repo repository.WorkoutRepository, jwtWrapper *utils.JwtWrapper) {
	workoutService := services.NewWorkoutService(repo, jwtWrapper)
	router.Get("/ping", workoutService.Ping)
	router.Post("/", workoutService.Create)
	router.Get("/", workoutService.GetWorkouts)
	router.Delete("/", workoutService.Delete)
}

func (a *WorkoutApp) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3002",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
