package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/t0nyandre/go-boilerplate-oauth2/pkg/database/postgres"
	"github.com/t0nyandre/go-boilerplate-oauth2/pkg/logger"
	"github.com/t0nyandre/go-boilerplate-oauth2/pkg/oauth2"
	"github.com/t0nyandre/go-boilerplate-oauth2/pkg/oauth2/github"
	"github.com/t0nyandre/go-boilerplate-oauth2/pkg/user"
)

func init() {
	if err := godotenv.Load("config/env/.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	ctx := context.Background()
	router := chi.NewRouter()
	logger := logger.NewLogger()
	postgres, err := postgres.NewPostgres(logger)
	if err != nil {
		logger.Fatalw("Failed to connect to database", "database", os.Getenv("POSTGRES_DB"), "error", err)
	}

	// Add to context
	ctx = context.WithValue(ctx, "logger", logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s", "World")
	})

	github := github.New(ctx, os.Getenv("OAUTH2_GITHUB_CLIENT_ID"), os.Getenv("OAUTH2_GITHUB_CLIENT_SECRET"), os.Getenv("OAUTH2_GITHUB_CALLBACK_URL"), "user")

	router.Mount("/auth", oauth2.NewRoutes(
		oauth2.NewService(oauth2.NewRepository(postgres), logger), logger))
	router.Mount("/user", user.NewRoutes(
		user.NewService(user.NewRepository(postgres), logger), logger))
	router.Mount("/auth/github", github.NewRoutes())
	logger.Infow("Successfully added routes", "auth", "/auth", "github", "/auth/github", "user", "/user")

	logger.Infow("Server successfully up and running", "host", os.Getenv("APP_HOST"), "port", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")), router); err != nil {
		logger.Fatalw("Server failed to start", "error", err)
	}
}