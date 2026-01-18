package walletservice

import (
	"context"

	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/google/uuid"
)

func (w *WalletService) GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	return nil, nil
}