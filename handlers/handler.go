package handler

import (
	"net/http"
	"strings"

	"github.com/atul-007/GreedyGameAssignment/models"
	"github.com/atul-007/GreedyGameAssignment/services"
	"github.com/gofiber/fiber/v2"
)

type DbHandlerInterface interface {
	HandleDb(ctx *fiber.Ctx) error
}
type DbHandler struct {
	Dbservices services.DbServicesInterface
}

func (d *DbHandler) HandleDb(ctx *fiber.Ctx) error {
	var req models.CommandRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid request"})
	}

	cmd := strings.Fields(req.Command)
	if len(cmd) == 0 {
		return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command"})
	}

	switch cmd[0] {
	case "SET":
		if len(cmd) < 3 {
			return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command format"})
		}

		expiryTime := d.Dbservices.ParseExpiryTime(cmd)
		condition := d.Dbservices.ParseCondition(cmd)

		err := d.Dbservices.Set(cmd[1], cmd[2], expiryTime, condition)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Error setting value"})
		}

		return ctx.JSON(models.CommandResponse{})

	case "GET":
		if len(cmd) < 2 {
			return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command format"})
		}

		value, err := d.Dbservices.Get(cmd[1])
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Error getting value"})
		}

		if value == "" {
			return ctx.Status(http.StatusNotFound).JSON(models.ErrorResponse{Error: "Key not found"})
		}

		return ctx.JSON(models.CommandResponse{Value: value})
	case "QPUSH":
		if len(cmd) < 3 {
			return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command format"})
		}

		values := cmd[2:]
		err := d.Dbservices.QPush(cmd[1], values)
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Error pushing values to queue"})
		}

		return ctx.JSON(models.CommandResponse{})
	case "QPOP":
		if len(cmd) < 2 {
			return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command format"})
		}

		value, err := d.Dbservices.QPop(cmd[1])
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Error popping value from queue"})
		}

		if value == "" {
			return ctx.Status(http.StatusNotFound).JSON(models.ErrorResponse{Error: "Queue is empty"})
		}
		return ctx.JSON(models.CommandResponse{Value: value})

	case "BQPOP":
		if len(cmd) < 2 {
			return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command format"})
		}
		value, err := d.Dbservices.BQPop(cmd[1], cmd[2])

		if value == "BUSY" {
			return ctx.Status(http.StatusServiceUnavailable).JSON(models.ErrorResponse{Error: "Queue is blocked"})
		}
		services.Check.Store("b")
		if err != nil {
			return ctx.Status(http.StatusInternalServerError).JSON(models.ErrorResponse{Error: "Error popping value from queue"})
		}

		if value == "" {
			return ctx.Status(http.StatusNotFound).JSON(models.ErrorResponse{Error: "Queue is empty"})
		}

		return ctx.JSON(models.CommandResponse{Value: value})

	default:
		return ctx.Status(http.StatusBadRequest).JSON(models.ErrorResponse{Error: "Invalid command"})
	}
}
