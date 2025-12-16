package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/rdssj/golang-assessment/internal/logger"
	"github.com/rdssj/golang-assessment/internal/middleware"
	"github.com/rdssj/golang-assessment/internal/routes"
)

func main() {
	lz := logger.New()
	defer lz.Sync()

	// Load .env if present for local development
	_ = godotenv.Load()

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/golang_assessment?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		lz.Fatal("failed to open db", zap.Error(err))
	}
	defer db.Close()

	if err := db.PingContext(context.Background()); err != nil {
		lz.Fatal("failed to ping db", zap.Error(err))
	}

	app := fiber.New()
	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(middleware.RequestLogger(lz))

	routes.Register(app, db, lz)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := app.Listen(":" + port); err != nil {
		lz.Fatal("fiber listen error", zap.Error(err))
	}
}
