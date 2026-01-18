package walletservice

import (
	"context"

	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
)

//go:generate mockery --name walletStorage
type WalletStorage interface {
	GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	ApplyOperation(ctx context.Context, op *models.WalletOperation) (error)
}

type WalletService struct {
	walletStorage WalletStorage
}

func NewWalletService(ctx context.Context, ws WalletStorage) (*WalletService) {
	return &WalletService{
		walletStorage: ws,
	}
}