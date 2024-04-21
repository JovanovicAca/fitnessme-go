package cmd

import (
	"context"
	"fitnessme/cors"
	"fitnessme/exercise/pkg/db"
	"fitnessme/exercise/pkg/repository"
	"fitnessme/exercise/pkg/services"
	"fmt"
	"net/http"

	"fitnessme/utils"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

type ExerciseApp struct {
	router http.Handler
	db.Handler
	repo repository.ExerciseRepository
	jw   utils.JwtWrapper
}

func New(jwtWrapper *utils.JwtWrapper) *ExerciseApp {
	handler := db.Init("localhost", "exercises", "postgres", "123", 5432)
	repo := repository.NewExerciseRepository(handler)

	app := &ExerciseApp{
		router: loadRoutes(repo, jwtWrapper),
		repo:   repo,
		jw:     *jwtWrapper,
	}

	app.router = cors.Cors(app.router)
	return app
}

func loadRoutes(repo repository.ExerciseRepository, jwtWrapper *utils.JwtWrapper) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/exercise", func(router chi.Router) {
		loadAuthRoutes(router, repo, jwtWrapper)
	})

	return router
}

func loadAuthRoutes(router chi.Router, repo repository.ExerciseRepository, jwtWrapper *utils.JwtWrapper) {
	exerciseService := services.NewExerciseService(repo, jwtWrapper)
	router.Get("/ping", exerciseService.Ping)
	router.Get("/byGroup", exerciseService.GetExerciseByGroupId)
	router.Post("/", exerciseService.Create)
	router.Post("/group", exerciseService.CreateGroup)
	router.Get("/group", exerciseService.GetAllGroups)
	router.Get("/", exerciseService.GetNameById)
	router.Get("/all", exerciseService.GetAllExercises)
	router.Delete("/", exerciseService.DeleteExercise)
	router.Put("/", exerciseService.Update)
}

func (a *ExerciseApp) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":3001",
		Handler: a.router,
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}
