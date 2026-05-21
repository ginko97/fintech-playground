package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ginko97/fintech-playground/internal/application"
	"github.com/ginko97/fintech-playground/internal/handler"
	"github.com/ginko97/fintech-playground/internal/infrastructure"
	"github.com/ginko97/fintech-playground/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// === Config ===
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://fintech:fintech123@localhost:5433/fintech_ledger?sslmode=disable"
	}

	// === PostgreSQL ===
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// === Redis ===
	redisClient, err := infrastructure.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// === Layers ===
	repo := repository.NewTransactionRepository(db)
	stateMachine := application.NewTransactionStateMachine(redisClient)
	workerPool := application.NewWorkerPool(5, stateMachine)
	workerPool.Start()
	defer workerPool.Shutdown()

	txService := application.NewTransactionService(repo, stateMachine, workerPool)

	txHandler := handler.NewTransactionHandler(txService)

	// === Router ===
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "transaction"})
	})

	v1 := r.Group("/api/v1")
	v1.POST("/transactions", txHandler.Create)

	// Graceful shutdown
	srv := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
