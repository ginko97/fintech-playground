package wallet

import (
	"context"
	"errors"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/ginko97/fintech-playground/internal/repository"
)

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

// CreateWallet creates a new wallet
func (s *WalletService) CreateWallet(ctx context.Context, userID, currency string) (*domain.Wallet, error) {
	wallet := domain.NewWallet(userID, currency)

	if err := s.repo.Create(ctx, wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// GetBalance returns current wallet balance with version
func (s *WalletService) GetBalance(ctx context.Context, walletID string) (*domain.Wallet, error) {
	return s.repo.FindByID(ctx, walletID)
}

// Credit adds money to wallet (with optimistic locking)
func (s *WalletService) Credit(ctx context.Context, walletID string, amount int64) error {
	wallet, err := s.repo.FindByID(ctx, walletID)
	if err != nil {
		return err
	}

	// Simple fraud check
	if amount > 10000000 { // 100 million IDR limit example
		return errors.New("amount exceeds single credit limit")
	}

	return s.repo.UpdateBalance(ctx, walletID, amount, wallet.Version)
}
