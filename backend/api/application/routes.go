package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/DanVerh/artschool-admin/backend/api/handlers"
)

// Create router with confgiured routes
func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/schedule", loadScheduleRoutes)
	router.Route("/students", loadStudentRoutes)

	return router
}

// Define all routes with HTTP methods
func loadStudentRoutes(router chi.Router) {
	studentHandler := &handler.StudentHandler{}
	router.Post("/", studentHandler.Create)
	router.Get("/", studentHandler.List)
	router.Get("/{id}", studentHandler.GetByID)
	router.Put("/{id}", studentHandler.UpdateByID)
	router.Delete("/{id}", studentHandler.DeleteByID)
}

func loadScheduleRoutes(router chi.Router) {
	scheduleHandler := &handler.Schedule{}
	router.Post("/", scheduleHandler.Create)
	router.Get("/", scheduleHandler.List)
	router.Get("/{id}", scheduleHandler.GetByID)
	router.Put("/{id}", scheduleHandler.UpdateByID)
	router.Delete("/{id}", scheduleHandler.DeleteByID)
}
