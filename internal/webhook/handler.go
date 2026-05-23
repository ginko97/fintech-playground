package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	secret string // webhook signing secret (from PSP)
}

func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{secret: secret}
}

// VerifySignature verifies PSP webhook signature (e.g. Xendit, Stripe)
func (h *WebhookHandler) VerifySignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

// HandleWebhook processes incoming webhook from external PSP
func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	_, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}

	signature := c.GetHeader("X-Signature") // example header
	if !h.VerifySignature(body, signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
		return
	}

	// TODO: Parse payload and update transaction state via state machine
	c.JSON(http.StatusOK, gin.H{
		"status":  "received",
		"message": "webhook processed successfully",
	})
}
