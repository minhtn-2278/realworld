package main

import (
	"context"

	"github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	echoSwagger "github.com/swaggo/echo-swagger/v2"

	"realworldapp/config"
	_ "realworldapp/docs"
	"realworldapp/internal/cache"
	"realworldapp/internal/handlers"
	appmiddleware "realworldapp/internal/middleware"
	"realworldapp/internal/repositories"
	"realworldapp/internal/services"
	database "realworldapp/pkg/db"
)

// @title Realworld API
// @version 1.0
// @description REST API for the Realworld application.
// @BasePath /api
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter `Bearer {access token}`.
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}
	redisStore, err := cache.NewRedisStore(context.Background(), cfg.RedisURL)
	if err != nil {
		panic(err)
	}
	defer redisStore.Close()

	e := echo.New()
	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Repositories
	articleRepository := repositories.NewArticleRepository(db)
	commentRepository := repositories.NewCommentRepository(db)
	tagRepository := repositories.NewTagRepository(db)
	userRepository := repositories.NewUserRepository(db)

	// Services
	articleService := services.NewArticleService(articleRepository, commentRepository, userRepository, db, redisStore)
	tagService := services.NewTagService(tagRepository, redisStore)
	userService := services.NewUserService(userRepository, cfg.JWTSecret)

	// Middleware
	authMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(cfg.JWTSecret),
		SigningMethod: echojwt.AlgorithmHS256,
		ErrorHandler:  appmiddleware.JWTErrorHandler,
	})

	// Handlers
	handlers.RegisterRoutes(e, articleService, tagService, userService, authMiddleware)

	if err := e.Start(cfg.HTTPAddr); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
