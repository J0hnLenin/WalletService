package walletserviceapi

import (
	"context"
	"errors"

	"github.com/J0hnLenin/WalletService/internal/models"

	custom_errors "github.com/J0hnLenin/WalletService/internal/errors"
	proto_models "github.com/J0hnLenin/WalletService/internal/pb/models"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (api *WalletServiceAPI) ApplyOperation(ctx context.Context, op *proto_models.Operation) (*proto_models.Wallet, error) {
	uuid, err := uuid.Parse(op.WalletId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "can't parse walletid '%s' request parameter %w", op.WalletId, err)
	}

	var amount int64
	
	if op.Amount < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount request parameter must be >= 0")
	}

	switch op.OperationType {
		case proto_models.Operation_WITHDRAW:
			amount = -op.Amount
		case proto_models.Operation_DEPOSIT:
			amount = op.Amount
		default:
			return nil, status.Errorf(codes.InvalidArgument, "operationType parameter must be '%s' or '%s', got '%s'", proto_models.Operation_DEPOSIT, proto_models.Operation_WITHDRAW, op.OperationType)
	}

	operation := &models.WalletOperation{
		WalletID: uuid,
		AmountChange: amount,
	}
	
	wallet, err := api.walletService.ApplyOperation(ctx, operation)
	if err != nil {
		if errors.Is(err, &custom_errors.ErrInsufficientBalance{}) {
			return nil, status.Error(codes.Aborted, err.Error())
		}
		return nil, err
	}

	return &proto_models.Wallet{
		Id:      wallet.ID.String(),
		Balance: wallet.Balance,
	}, nil
}