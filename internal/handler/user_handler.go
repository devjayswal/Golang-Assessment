package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/rdssj/golang-assessment/internal/models"
	"github.com/rdssj/golang-assessment/internal/service"
)

type UserHandler struct {
	svc *service.UserService
	lz  *zap.Logger
}

func NewUserHandler(svc *service.UserService, lz *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, lz: lz}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
		Dob  string `json:"dob"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	d, err := time.Parse("2006-01-02", body.Dob)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dob must be YYYY-MM-DD")
	}
	user, err := h.svc.Create(context.Background(), models.User{Name: body.Name, Dob: d})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(models.UserResponse{ID: user.ID, Name: user.Name, Dob: user.Dob.Format("2006-01-02")})
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	resp, err := h.svc.Get(context.Background(), id)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return c.JSON(resp)
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	var body struct {
		Name string `json:"name"`
		Dob  string `json:"dob"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid json")
	}
	d, err := time.Parse("2006-01-02", body.Dob)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "dob must be YYYY-MM-DD")
	}
	user, err := h.svc.Update(context.Background(), id, models.User{Name: body.Name, Dob: d})
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(models.UserResponse{ID: user.ID, Name: user.Name, Dob: user.Dob.Format("2006-01-02")})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(context.Background(), id); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "user not found")
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	limit := int32(c.QueryInt("limit", 50))
	offset := int32(c.QueryInt("offset", 0))
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	items, err := h.svc.List(context.Background(), limit, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(items)
}
