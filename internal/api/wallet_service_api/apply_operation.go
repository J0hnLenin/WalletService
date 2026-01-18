package walletserviceapi

import (
	"context"

	"github.com/J0hnLenin/WalletService/internal/models"
	proto_models "github.com/J0hnLenin/WalletService/internal/pb/models"
	"github.com/google/uuid"
)

func (api *WalletServiceAPI) ApplyOperation(ctx context.Context, op *proto_models.Operation) (*proto_models.Wallet, error) {
	uuid := uuid.UUID([]byte(op.WalletId))

	amount := op.Amount
	if op.OperationType == proto_models.Operation_WITHDRAW {
		amount = -amount
	}

	operation := &models.WalletOperation{
		WalletID: uuid,
		AmountChange: amount,
	}
	
	wallet, err := api.walletService.ApplyOperation(ctx, operation)
	if err != nil {
		return nil, err
	}

	return &proto_models.Wallet{
		Id:      wallet.ID.String(),
		Balance: wallet.Balance,
	}, nil
}