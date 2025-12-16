package routes

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	db "github.com/rdssj/golang-assessment/db/sqlc"
	"github.com/rdssj/golang-assessment/internal/handler"
	"github.com/rdssj/golang-assessment/internal/service"
)

func Register(app *fiber.App, dbConn *sql.DB, lz *zap.Logger) {
	repo := db.New(dbConn)
	svc := service.NewUserService(repo)
	h := handler.NewUserHandler(svc, lz)

	app.Post("/users", h.Create)
	app.Get("/users/:id", h.Get)
	app.Put("/users/:id", h.Update)
	app.Delete("/users/:id", h.Delete)
	app.Get("/users", h.List)
}
