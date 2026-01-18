package bootstrap

import (
	server "github.com/J0hnLenin/WalletService/internal/api/wallet_service_api"
	walletservice "github.com/J0hnLenin/WalletService/internal/services/wallet_service"
)

func InitWalletServiceAPI(walletService *walletservice.WalletService) *server.WalletServiceAPI {
	return server.NewWalletServiceAPI(walletService)
}