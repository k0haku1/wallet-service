package middleware

import (
	"github.com/gofiber/fiber/v2"
	"log"
	svcErrors "wallet-service/internal/wallet/service/errors"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	switch err {
	case svcErrors.ErrWalletNotFound:
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	case svcErrors.ErrInvalidOperation:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	case svcErrors.ErrInsufficientFunds:
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case svcErrors.ErrInvalidAmount:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:

		log.Printf("unexpected error: %v", err)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
}
