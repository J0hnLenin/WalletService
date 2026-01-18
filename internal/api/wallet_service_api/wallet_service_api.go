package walletserviceapi

import (
	"context"

	"github.com/J0hnLenin/WalletService/internal/models"
	"github.com/J0hnLenin/WalletService/internal/pb/wallets_api"
	"github.com/google/uuid"
)

type walletService interface {
	GetWalletByID(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
	ApplyOperation(ctx context.Context, op *models.WalletOperation) (*models.Wallet, error)
}

type WalletServiceAPI struct {
	wallets_api.UnimplementedWalletsServiceServer
	walletService walletService
}

func NewWalletServiceAPI(ws walletService) *WalletServiceAPI {
	return &WalletServiceAPI{
		walletService: ws,
	}
}
