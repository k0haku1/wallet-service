package handler

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
	"wallet-service/internal/dto"
	"wallet-service/internal/wallet/service"
)

type WalletHandler struct {
	svc *service.WalletService
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{
		svc: svc,
	}
}

func (h *WalletHandler) GetWalletBalance(c *fiber.Ctx) error {
	walletId, err := uuid.Parse(c.Params("wallet_uuid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid wallet UUID"})
	}

	balance, err := h.svc.GetWallet(c.Context(), walletId)
	if err != nil {
		return err
	}

	return c.JSON(dto.GetWalletResponse{
		WalletID: walletId,
		Balance:  balance,
	})
}

func (h *WalletHandler) UpdateWalletBalance(c *fiber.Ctx) error {
	var req dto.WalletOperationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	op, err := h.svc.UpdateWalletBalance(ctx, req)
	if err != nil {
		return err
	}

	resp := dto.WalletOperationResponse{
		OperationID:   op.ID,
		WalletID:      op.WalletID,
		OperationType: op.Type,
		Amount:        op.Amount,
		CreatedAt:     op.CreatedAt.Format(time.RFC3339),
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
