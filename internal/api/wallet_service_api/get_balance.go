package walletserviceapi

import (
	"context"

	proto_models "github.com/J0hnLenin/WalletService/internal/pb/models"
	"github.com/J0hnLenin/WalletService/internal/pb/wallets_api"
	"github.com/google/uuid"
)

func (api *WalletServiceAPI) GetBalance(ctx context.Context, req *wallets_api.GetBalanceRequest) (*proto_models.Wallet, error){
	uuid := uuid.UUID([]byte(req.WalletId))
	wallet, err := api.walletService.GetWalletByID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	
	return &proto_models.Wallet{
		Id: wallet.ID.String(),
		Balance: wallet.Balance,
	}, nil
}