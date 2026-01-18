package bootstrap

import (
	"context"

	walletservice "github.com/J0hnLenin/WalletService/internal/services/wallet_service"
	"github.com/J0hnLenin/WalletService/internal/storage/pgstorage"
)

func InitWalletService(storage *pgstorage.PGStorage) *walletservice.WalletService {
	return walletservice.NewWalletService(context.Background(), storage)
}