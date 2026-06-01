package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground/internal/application"
	"github.com/ginko97/fintech-playground/internal/fraud"
	"github.com/ginko97/fintech-playground/internal/handler"
	"github.com/ginko97/fintech-playground/internal/infrastructure"
	"github.com/ginko97/fintech-playground/internal/middleware"
	"github.com/ginko97/fintech-playground/internal/repository"
	"github.com/ginko97/fintech-playground/internal/wallet"
	"github.com/ginko97/fintech-playground/internal/webhook"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// === Initialize Logger & Config ===
	infrastructure.InitLogger()
	defer infrastructure.SyncLogger()

	logger := infrastructure.GetLogger()
	logger.Info("Starting fintech playground service...")

	config, err := infrastructure.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	// === Database ===
	db, err := pgxpool.New(context.Background(), config.Database.URL)
	if err != nil {
		logger.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}
	defer db.Close()

	// === Redis ===
	redisClient, err := infrastructure.NewRedisClient()
	if err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
	}
	defer redisClient.Close()

	// === Layers ===
	repo := repository.NewTransactionRepository(db)

	// Wallet
	walletRepo := repository.NewWalletRepository(db) // ← Add this
	walletService := wallet.NewWalletService(walletRepo)

	// Fraud
	fraudEngine := fraud.NewFraudEngine()

	stateMachine := application.NewTransactionStateMachine(redisClient)
	workerPool := application.NewWorkerPool(5, stateMachine)
	workerPool.Start()
	defer workerPool.Shutdown()

	txService := application.NewTransactionService(repo, stateMachine, workerPool, fraudEngine, walletService)
	txHandler := handler.NewTransactionHandler(txService)
	webhookHandler := webhook.NewWebhookHandler("your-webhook-secret-123")

	// === Router ===
	r := gin.Default()

	if config.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Middlewares
	r.Use(middleware.RequestIDMiddleware()) // ← Add this
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.RateLimit(redisClient, 100, 1*time.Minute))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/api/v1")
	v1.Use(middleware.RateLimit(redisClient, 100, 1*time.Minute))

	v1.POST("/transactions", txHandler.Create)
	v1.POST("/webhooks/psp", webhookHandler.HandleWebhook)

	// === Start Server ===
	srv := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	logger.Info("Server started", zap.String("port", config.Server.Port))

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}
}
