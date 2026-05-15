package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground/internal/application"
)

type TransactionHandler struct {
	service *application.TransactionService
}

func NewTransactionHandler(service *application.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

type CreateTransactionResponse struct {
	TransactionID  string `json:"transaction_id"`
	Status         string `json:"status"`
	IdempotencyKey string `json:"idempotency_key"`
}

func (h *TransactionHandler) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req application.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation (can move to validator later)
	if req.IdempotencyKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idempotency_key is required"})
		return
	}

	tx, err := h.service.Create(ctx, req)
	if err != nil {
		// TODO: proper error mapping (conflict, internal, etc)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, CreateTransactionResponse{
		TransactionID:  tx.ID,
		Status:         string(tx.Status),
		IdempotencyKey: tx.IdempotencyKey,
	})
}
