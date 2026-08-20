package main

import (
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"

	"realworldapp/config"
	"realworldapp/internal/handlers"
	database "realworldapp/pkg/db"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// Keep the database connection available for repositories and handlers.
	_ = db

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	handlers.RegisterRoutes(e)

	if err := e.Start(cfg.HTTPAddr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
