package application

import (
  "context"
  "fmt"
  "net/http"
  "strconv"
)

// define port constant value
const port int = 8080

type App struct {
  router http.Handler
}

func New() *App {
  app := &App{
    router: loadRoutes(),
  }
  
  return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port), // convert port to ASCII
		Handler: a.router,
	}

	err := server.ListenAndServe()
	if err != nil {
    return fmt.Errorf("failed to start server: %w", err)
	}
  
  return nil
}