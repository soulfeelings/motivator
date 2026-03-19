package response

import (
	"github.com/gofiber/fiber/v3"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Meta    Meta   `json:"meta"`
	Error   string `json:"error,omitempty"`
}

type Meta struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func Success(c fiber.Ctx, data any) error {
	return c.JSON(Response{Success: true, Data: data})
}

func Created(c fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Success: true, Data: data})
}

func Paginated(c fiber.Ctx, data any, meta Meta) error {
	return c.JSON(PaginatedResponse{Success: true, Data: data, Meta: meta})
}

func Error(c fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(Response{Success: false, Error: msg})
}

func BadRequest(c fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusBadRequest, msg)
}

func Unauthorized(c fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusUnauthorized, msg)
}

func Forbidden(c fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusForbidden, msg)
}

func NotFound(c fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusNotFound, msg)
}

func Conflict(c fiber.Ctx, msg string) error {
	return Error(c, fiber.StatusConflict, msg)
}

func InternalError(c fiber.Ctx) error {
	return Error(c, fiber.StatusInternalServerError, "internal server error")
}
